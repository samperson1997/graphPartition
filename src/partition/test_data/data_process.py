f = open("comljungraph.txt")


mp = dict()
id = 0
cnt = 0
for line in f.readlines():
    line = line.strip()
    line = line.split()
    if cnt == 0:
        vertex = int(line[0])
    else:
        src,dst = int(line[0]),int(line[1])
        srcc = 0
        if mp.__contains__(src):
            srcc = mp[src]
        else:
            srcc = id
            mp[src] = id
            id+=1 
        dstt = 0
        if mp.__contains__(dst):
            dstt = mp[dst]
        else:
            dstt = id
            mp[dst] = id
            id+=1

        print(srcc,dstt)
    cnt+=1
print(len(mp))