array = []
threads = []
4.times do
    Thread.new do
        2500.times do |n| 
            array << n
        end
    end
end
Join
puts array.size
