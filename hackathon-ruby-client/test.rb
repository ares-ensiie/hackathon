require 'ares-hackathon'

c = AresHackathon.new
c.connect("127.0.0.1", "Salut")
while c.status == AresHackathon::ONGOING
  puts "--- Waiting for my turn"
  c.next_turn

  break if c.status != AresHackathon::ONGOING

  puts "--- Attacking"
  m = c.my_cells
  for i in 0..10
    cell = m.sample
    dx = cell.x + rand(3) - 1
    dy = cell.y + rand(3) - 1
    if dx >= 0 && dx < c.field.size_x && dy >= 0 && dy < c.field.size_y then
      c.attack cell.x, cell.y, dx, dy
    end
  end
  puts "--- Ending attacks"
  c.end_attacks

  puts "--- Adding units"
  while c.units_remained > 0
    c.add_units c.my_cells.sample, 1
  end
  c.end_adding_units
end

puts "--- WIN" if c.status == AresHackathon::VICTORY
puts "--- LOST" if c.status == AresHackathon::DEFEAT
puts "--- NETWORK" if c.status == AresHackathon::CONNECTION_LOST
puts "--- WTF" if c.status == AresHackathon::ONGOING
