require 'socket'

class AresHackathon
  ONGOING         = 0
  VICTORY         = 1
  DEFEAT          = 2
  CONNECTION_LOST = 3

  attr_reader :field
  attr_reader :status
  attr_reader :name
  attr_reader :id

  def initialize
    @status = ONGOING
  end

  def connect(url, name)
    s = url.split ":"
    ip = s[0]
    port = s.length > 1 ? s[1].to_i : 1337

    @conn = TCPSocket.new ip, port
    name = name[0..24]
    name = name+"\0"

    @conn.send name, 0
    id = @conn.recv(1)

    @id = id.ord
    @name = name[0..24]
  end


  def attack(fromX, fromY, toX, toY)
    buffer = [fromX, fromY, toX, toY].pack('C*')
    @conn.send buffer, 0

    if fromX == 255 && fromY == 255 && toX == 255 && toY == 255 then
      return
    end

    @field = AresHackathon::Field.read_from_socket @conn
  end

  def end_attacks
    attack(255,255,255,255)

    buffer = @conn.recv 1
    @total_units = buffer[0].ord
    @units = []
  end

  def units_remained
    return @total_units - @units.length
  end

  def add_unit(c)
    if units_remained > 0 then
      @units << c
    end
  end

  def add_units_list(cells)
    cells.each do | cell |
      add_unit cell
    end
  end

  def add_units(c, count)
    for i in 1..count
      add_unit c
    end
  end

  def end_adding_units
    buffer = []
    for i in 0..(@total_units-1)
      if i < @units.length then
        buffer << @units[i].x
        buffer << @units[i].y
      else
        buffer << 255
        buffer << 255
      end
    end
    send_buffer = buffer.pack("C*")
    @conn.write send_buffer

    @field = AresHackathon::Field.read_from_socket @conn
  end

  def next_turn
    @field = AresHackathon::Field.read_from_socket @conn

    @status = VICTORY if @field.has_win @id
    @status = DEFEAT if @field.has_lost @id
  end

  def cell(x,y)
    return @field.cell(x,y)
  end

  def my_cells
    return @field.owned_by @id
  end

  def map
    return @field
  end
end

require 'ares-hackathon/field'
require 'ares-hackathon/cell'
