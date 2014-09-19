r = 1000000
rsqr = r*r
iteration = 10000

cnt = [0, 0]
16.times do
    Thread.new do
        iteration.times do 
            x = rand(2*r) - r
            x = x*x
            y = rand(2*r) - r
            y = y*y
            i = (x+y) / rsqr 
            cnt.at(i).atomic_add(1)
        end
    end
end
Join
num = cnt[0] * 4
den = iteration * 16

puts num / den.to_f
