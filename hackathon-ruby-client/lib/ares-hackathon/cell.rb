class AresHackathon::Cell
  attr_reader :x
  attr_reader :y
  attr_reader :owner
  attr_reader :power
  
  def initialize(px,py, owner, power)
    @x = px
    @y = py
    @owner = owner
    @power = power
  end
end
