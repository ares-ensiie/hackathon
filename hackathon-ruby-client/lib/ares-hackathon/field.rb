class AresHackathon::Field
  attr_reader :size_x
  attr_reader :size_y
  def initialize(sx, sy, field)
    @size_x = sx
    @size_y = sy
    @field  = field
  end

  def self.read_from_socket(conn)
    sizes = conn.recv(2)
    sx = sizes[0].ord
    sy = sizes[1].ord
    field = []
    for x in 0..(sx-1)
      field[x] = []
    end

    for y in 0..(sy-1)
      buffer = conn.recv(2*sx)
      for x in 0..(sx-1)
        field[x][y] = AresHackathon::Cell.new x, y, buffer[2*x].ord, buffer[2*x+1].ord
      end
    end

    return AresHackathon::Field.new sx, sy, field
  end

  def cell(x,y)
    return @field[x][y]
  end

  def has_lost(p)
    @field.each do |row|
      row.each do |cell|
        return false if cell.owner == p
      end
    end
    return true
  end

  def has_win(p)
    @field.each do |row|
      row.each do |cell|
        return false if cell.owner != p && cell.owner != 0
      end
    end
    return true
  end

  def owned_by(player)
    res =  @field.map do |line|
      line.select do |cell|
        cell.owner == player
      end
    end

    return res.flatten
  end
end
