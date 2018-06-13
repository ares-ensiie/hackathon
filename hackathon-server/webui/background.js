(function() {
	// strict mode
	"use strict";
	// configurable options that determine particle behavior and appearance
	var options = {
			// connector line options
			connector: {
				// color of line
				color: "white",
				// base line width of each connector
				lineWidth: 0.5
			},
			// individual particle options
			particle: {
				// base color
				color: "white",
				// number of particles to render onto stage
				count: 75,
				// size of particle influence when drawing connectors (in pixels)
				influence: 75,
				// a range of sizes for an individual particle (falls within max-min)
				sizeRange: {
					min: 1,
					max: 5
				},
				// a range of velocities for each particle to inherit (falls within max-min)
				velocityRange: {
					min: 10,
					max: 60
				}
			}
		},
		particles,
		stage;
	// sets the render loop state based on page's visibility
	function handleDocumentVisibilityChange() {
		document.hidden ?
			stage.render.pause() :
			stage.render.resume();
	}

	function initPage() {
		initObjects();
		initListeners();
	}
	// invoked when page has loaded
	// registers window and document level event listeners
	// required post object instantiation
	function initListeners() {
		// resize canvas on window resize
		window.addEventListener(
			"resize",
			stage.resize);
		// pause and resume rendering loop when visibility state is changed
		document.addEventListener(
			"visibilitychange",
			handleDocumentVisibilityChange);
	}
	// invoked when page has loaded
	// performs setup of stage and particles
	// triggers the rendering loop
	function initObjects() {
		particles = new ParticleGroup();
		// create our stage instance
		stage = new Stage(document.querySelector("canvas"));
		// init the canvas bounds and fidelity
		stage.resize();
		// populate particle group collection
		particles.populate();
		// begin stage rendering with the renderAll draw routine on each tick
		stage.render.begin(
			particles.render.bind(particles));
	}
	// renders an FPS value to the canvas top-left corner
	function renderFPS(value) {
		context.save();
		context.font = 10 * stage.pixelRatio + "px sans-serif";
		context.fillStyle = "white";
		context.textBaseline = "top";
		context.fillText(value, 5 * stage.pixelRatio, 5 * stage.pixelRatio);
		context.restore();
	}
	// Paricle object construct:
	// stores critical particle information regarding trajectory and location
	// 		v: distance delta coefficient
	// 		t: angular theta in degrees
	//		r: particle radius
	//		p: position on stage, x:y coordinate properties
	// 		influence: radial influence of connectivity to other particles
	// defines methods relating to particle rendering and comparison to other particle locations
	function Particle(props) {
		// private particle properties
		var x = props.x || 0,
			y = props.y || 0,
			v = props.speed,
			t = props.theta,
			r = props.radius,
			influence = props.influence,
			color = props.color;
		// public properties and methods
		return {
			// position property getter: read-only
			get x() {
				return x;
			},
			get y() {
				return y;
			},
			// compares this particle location to another
			// returns a coefficient describing the ratio of closeness to particle influence
			influencedBy: function(a, b) {
				var hyp = Math.sqrt(Math.pow(a - x, 2) + Math.pow(b - y, 2));
				return Math.abs(hyp) <= influence ?
					hyp / influence : 0;
			},
			// render this particle to a given context
			render: function(ctx) {
				ctx.strokeStyle = color;
				ctx.lineWidth = 1 / r * r * stage.pixelRatio;
				ctx.beginPath();
				ctx.arc(x || 10, y || 10, r, 0, 8);
				ctx.stroke();
			},
			// sets the position of this particle compared to its previous position
			// in relation to a total time delta since last positioning
			// positions infinitely loop within canvas bounds
			setPosition: function(timeDelta) {
				x += timeDelta / v * Math.cos(t);
				x = x > stage.width + r ? 0 - r : x;
				x = x < 0 - r ? stage.width + r : x;
				y += timeDelta / v / 2 * Math.sin(t);
				y = y > stage.height + r ? 0 - r : y;
				y = y < 0 - r ? stage.height + r : y;
			}
		}
	}
	// ParticleGroup object construct
	// an object which controls higher level methods to create, destroy and compare
	// a grouping of Particle instances. Holds a private collection of all existing particles
	function ParticleGroup() {
		// particle instance array
		var _collection = [],
			// cache connector option property values
			_connectColor = options.connector.color,
			_connectWidth = options.connector.lineWidth,
			_particleCount = options.particle.count;
		// generates and returns a new particle instance with
		// randomized properties within ranges defined in the options object
		function _generateNewParticle() {
			return new Particle({
				x: stage.width * Math.random(),
				y: stage.height * Math.random(),
				speed: getRandomWithinRange(options.particle.velocityRange) / stage.pixelRatio,
				radius: getRandomWithinRange(options.particle.sizeRange),
				theta: Math.round(Math.random() * 360),
				influence: options.particle.influence * stage.pixelRatio,
				color: options.particle.color
			});
			// random number generator within a given range object defined by a max and min property
			function getRandomWithinRange(range) {
				return ((range.max - range.min) * Math.random() + range.min) * stage.pixelRatio;
			}
		}
		// loops through particle collection
		// sets the new particle location given a time delta
		// queries other particles to see if a connection between particles should be rendered
		function _checkForNeighboringParticles(ctx, p1) {
			// cache particle influence method
			var getInfluenceCoeff = p1.influencedBy;
			// particle collection iterator
			for (var i = 0, p2, d; i < _particleCount; i++) {
				p2 = _collection[i];
				// conditional checks - ignore if p2 is the same as p1
				// or if p2 has already been iterated through to check for neighbors
				if (p1 !== p2 && !p2.checked) {
					// compare the distance delta between the two particles
					d = getInfluenceCoeff(p2.x, p2.y);
					// render the connector if coefficient is non-zero
					if (d) {
						_connectParticles(ctx, p1.x, p1.y, p2.x, p2.y, d);
					}
				}
			}
		}
		// given two particles and an influence coefficient between them,
		// renders a connecting line between the two particles
		// the coefficient determines the opacity (strength) of the connection
		function _connectParticles(ctx, x1, y1, x2, y2, d) {
			ctx.save();
			ctx.globalAlpha = 1 - d;
			ctx.strokeStyle = _connectColor;
			ctx.lineWidth = _connectWidth * (1 - d) * stage.pixelRatio;
			ctx.beginPath();
			ctx.moveTo(x1, y1);
			ctx.lineTo(x2, y2);
			ctx.stroke();
			ctx.restore();
		}
		// public object
		return {
			// returns the size of the collection
			get length() {
				return _collection.length;
			},
			// adds a particle object instance to the collection
			add: function(p) {
				_collection.push(p || _generateNewParticle());
			},
			// initial population of bar collection
			populate: function() {
				for (var i = 0; i < _particleCount; i++) {
					this.add();
				}
			},
			// loops through all particle instances within collection and
			// invokes instance rendering method on a given context at tick coefficient t
			render: function(ctx, t) {
				// loop through each particle, position it and render
				for (var i = 0, p; i < _particleCount; i++) {
					p = _collection[i];
					p.checked = false;
					p.setPosition(t);
					p.render(ctx);
				}
				// loop through each particle, check for connectors to be rendered with neighboring particles
				for (var i = 0, p; i < _particleCount; i++) {
					p = _collection[i];
					_checkForNeighboringParticles(ctx, p);
					p.checked = true;
				}
			}
		};
	}
	// Stage object construct
	// provides a structure to store a canvas and its context upon isntantiation
	// provides methods to interact with and adjust the canvas
	function Stage(canvasEl) {
		// canvas element reference
		var canvas = canvasEl instanceof Node ?
			canvasEl :
			document.querySelector(canvasEl),
			// 2d context reference
			context = canvas.getContext("2d"),
			// cache the device pixel ratio
			pixelRatio = window.devicePixelRatio;
		// resize the canvas initially to fit its containing element
		function _adjustCanvasBounds() {
			var parentSize = canvas.parentNode.getBoundingClientRect();
			canvas.width = parentSize.width;
			canvas.height = parentSize.height;
		}
		// updates the canvas dimensions relative to the device's pixel ratio for highest fidelity and accuracy
		function _adjustCanvasFidelity() {
			canvas.style.width = canvas.width + "px";
			canvas.style.height = canvas.height + "px";
			canvas.width *= pixelRatio;
			canvas.height *= pixelRatio;
		}
		// public object returned
		return {
			// dimension getters
			get height() {
				return canvas.height;
			},
			get width() {
				return canvas.width;
			},
			get pixelRatio() {
				return pixelRatio;
			},
			init: function(el) {
				canvas = el;
				context = canvas.getContext("2d");
			},
			resize: function() {
				_adjustCanvasBounds();
				_adjustCanvasFidelity();
			},
			render: (function() {
				// a flag indicating state of animation loop
				var paused = false,
					// animation request ID reported during loop
					requestID = null,
					// reference to a function detailing all draw operations per frame of animation
					renderMethod;
				// public methods
				return {
					// once invoked, creates a rendering loop which in turn invokes
					// drawMethod argument, passing along the delta (in milliseconds) since the last frame draw
					begin: function(drawMethod) {
						// cache the draw method to invoke per frame
						renderMethod = drawMethod;
						// a reference to requestAnimationFrame object
						// this should also utilize some vendor prefixing and a setTimeout fallback
						var requestFrame = window.requestAnimationFrame,
							// get initial start time
							latestTime, startTime = Date.now(),
							// during each interval, clear the canvas and invoke renderMethod
							intervalMethod = function(tick) {
								this.clear();
								renderMethod(
									context,
									tick);
							}.bind(this);
						// start animation loop
						(function loop() {
							// calculate tick time between frames
							var now = Date.now(),
								tick = now - latestTime || 1;
							// update latest time stamp
							latestTime = now;
							// report tick value to intervalCallback
							intervalMethod(tick);
							// loop iteration if no pause state is set
							requestID = paused ?
								null :
								requestAnimationFrame(loop);
						})();
					},
					// clears the stage's canvas
					clear: function() {
						context.clearRect(0, 0, canvas.width, canvas.height);
					},
					// pause the canvas rendering
					pause: function() {
						paused = true;
						cancelAnimationFrame(requestID);
					},
					// resumes the animation with a given rendering method
					resume: function() {
						paused = false;
						this.begin(renderMethod);
					}
				}
			}())
		};
	}
	// our init is triggered when the window is loaded/ready
	window.addEventListener(
		"load",
		initPage);
}());
