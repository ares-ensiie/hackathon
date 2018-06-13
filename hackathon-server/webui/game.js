var Config = {
  cell: {
    sizeX: 50,
    sizeY: 50,
    margin: 10
  },
  teamColors: [
    "white",
    "#00B8E0", // sky blue
    "#f48f30", // orange
    "#ce4e4e", // red
    "#f2f23a", // yellow
    "#c044ea", // purple
    "#417231", // dark green
    "#ef7ada", // pink
    "#11c1b3" // blue green
  ]
}

var Game;
var Events = [];

function getCell(x, y) {
  return $('#cell-'+x+'-'+y);
}

function getLevelup(x, y) {
  return $('#cell-'+x+'-'+y+' span.lvlup');
}

jQuery.fn.getPop = function() {
    var e = $(this[0]) // It's your element
    return parseInt(e.find('.pop').html());
};

jQuery.fn.setPop = function(pop) {
    var e = $(this[0]) // It's your element
    e.find('.pop').html(pop);
    return this; // This is needed so others can keep chaining off of this
};

function nextEvent() {
  var event = Events.shift();
  if (event === undefined) {
    setTimeout(function(){
      nextEvent();
    }, 100);
  } else {
    switch(event.type) {
      case "attack":
        if(event.data) {
          attack(event.data, function() {
          })
        }
        nextEvent();
        break;
      case "initplacement":
        if(event.data) {
          initPlacement(event.data);
        }
        nextEvent();
        break;
      case "placement":
        if(event.data) {
          placement(event.data);
        }
        nextEvent();
        break;
      case "log":
        if(event.data) {
          addLog(event.data);
        }
        nextEvent();
        break;
      case "field":
        if(event.data) {
          updateBoard(event.data.sizeX, event.data.sizeX, event.data.cells);
        }
        nextEvent();
        break;
      default:
        console.log("event not implemented");
        console.log(event.type);
    }
  }
}

function Event(type, data) {
  this.type = type;
  this.data = data;
}

function initBoard(sizeX, sizeY, values) {
  $('#board .row').remove();
  for (var y = 0; y < sizeY; y++) {
    $("#board").append('<div id="row-'+y+'" class="row"></div>');
    for (var x = 0; x < sizeX; x++) {
      $('#row-'+y).append('<div id="cell-'+x+'-'+y+'" class="cell"></div>');
      var htmlCell = getCell(x, y);
      var cell = values[x][y];
      htmlCell.data('owner', cell.owner.ID);
      htmlCell.css('background-color', Config.teamColors[cell.owner.ID]);
      htmlCell.append('<div class="pop">'+cell.population+'</div>');
    }
  }

  $('#board').width(sizeX*(Config.cell.sizeX+Config.cell.margin*2));
}

function updateBoard(sizeX, sizeY, values) {
  for (var y = 0; y < sizeY; y++) {
    for (var x = 0; x < sizeX; x++) {
      var cell = getCell(x, y);
      cell.setPop(values[x][y].population);
      cell.css('background-color', Config.teamColors[values[x][y].owner.ID]);
    }
  }
}

function initTeams(teams) {
  $('#teams ul li').remove();
  for (var i = 0; i < teams.length; i++) {
    addTeam(teams[i]);
  }
}

function addTeam(team) {
  if(team) {
    $('#teams ul').append('<li id="team-'+team.ID+'">'+team.name+'</li>');
    $('#team-'+team.ID).css('color', Config.teamColors[team.ID]);
  }
}

function addLog(log) {
  if(log) {
    var style = ""
    if(log.id != undefined){
      style +="color: "+Config.teamColors[log.id]+"; "
    }
    if(log.level == "error") {
      style += "font-weight: bold; text-decoration: underline; "
    }
    $('#logs ul').append('<li class="log" style="'+style+'">'+log.msg+'</li>');
    $("#logs").scrollTop($("#logs").prop("scrollHeight"));
  }
}

function setCurrentTeam(id) {
  $('#teams ul li').css('border-width', '0px');
  $('#team-'+id).css('border-width', '2px');
}

function initPlacement(placement) {
  var htmlCell = getCell(placement.x, placement.y);
  htmlCell.data('owner', placement.player.ID);
  htmlCell.css('background-color', Config.teamColors[placement.player.ID]);
  htmlCell.setPop(placement.pop);
}

function shake(x, y, callback) {
  var elem = document.getElementById('cell-'+x+'-'+y);
  elem.classList.remove("shake");
  void elem.offsetWidth; // trick
  getCell(x, y).addClass("shake").delay(820).queue(function(){
    if(callback) {
      callback();
    }
    $(this).dequeue();
  });

}

function attack(a, callback) {
  var fromCell = getCell(a.fromX, a.fromY);
  var fromValue = fromCell.getPop();
  var toCell = getCell(a.toX, a.toY);
  var toValue = toCell.getPop();

  setCurrentTeam(a.attacker.ID);
  shake(a.fromX, a.fromY);
  shake(a.toX, a.toY, function() {
    // if (a.isWon) {
    //   toCell.data('owner', a.attacker.ID);
    //   toCell.css('background-color', Config.teamColors[a.attacker.ID]);
    // }
    // fromCell.setPop(a.attackerPop);
    // toCell.setPop(a.defenderPop);
    lvlUp(a.fromX, a.fromY, (a.attackerPop-fromValue), false);
    var defenderDiff = a.defenderPop-toValue;
    if (defenderDiff < 0) {
      lvlUp(a.toX, a.toY, defenderDiff, false);
    } else if (defenderDiff > 0) {
      lvlUp(a.toX, a.toY, defenderDiff, true);
    }

    if (callback) {
      callback();
    }
  });
}

function lvlUp(x, y, text, isBonus) {
  if (isBonus) {
    getCell(x, y).append('<span class="lvlup buff">'+text+'</span>');
  } else {
    getCell(x, y).append('<span class="lvlup debuff">'+text+'</span>');
  }
}

function placement(p) {
  setCurrentTeam(p.player.ID);
  var cell = getCell(p.x, p.y);
  // cell.setPop(cell.getPop()+1);
  lvlUp(p.x, p.y, "up", true);
}

$( document ).ready(function() {
  var socket = io();
  nextEvent();

  socket.on('game init', function(msg){
    var game = JSON.parse(msg);
    Game = game;
    Events = [];
    initTeams(game.players);
    initBoard(game.field.sizeX, game.field.sizeX, game.field.cells);
  });

  // {"attacker":{"ID":1,"Name":"Slt9\u0000"},"fromX":0,"fromY":1,"toX":1,"toY":0,
  // "isWon":true,"attackerPop":1,"defenderPop":2}
  socket.on('game attack', function(msg){
    var a = JSON.parse(msg);
    Events.push(new Event("attack", a));
  });

  socket.on('game initplacement', function(msg){
    var p = JSON.parse(msg);
    addTeam(p.player);
    Events.push(new Event("initplacement", p));
  });

  socket.on('game placement', function(msg){
    var p = JSON.parse(msg);
    if(p == null) {
      console.log("PLACEMENT NULL");
    } else {
      Events.push(new Event("placement", p));
    }
  });

  socket.on('game log', function(msg){
    var log = JSON.parse(msg);
    Events.push(new Event("log", log));
  });

  socket.on('game field', function(msg){
    var field = JSON.parse(msg);
    console.log(field);
    Events.push(new Event("field", field));
  });
});
