import sys
'''
data_process.py [data_path]
'''
mp = dict()
if len(sys.argv)!=2:
    print("error no input file")
    exit(0)
id,cnt = 0,0
data_path = sys.argv[1]
f = open(data_path)
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