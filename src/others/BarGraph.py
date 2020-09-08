import matplotlib.pyplot as plt
import matplotlib
matplotlib.rcParams['font.sans-serif'] = ['SimHei']
matplotlib.rcParams['axes.unicode_minus'] = False

label_list = ['2', '5', '10']
num_list1 = [1.47, 2.24, 2.89]
num_list2 = [1.45, 2.12, 2.71]
num_list3 = [1.24, 1.69, 1.95]
x = range(len(num_list1))
x = [i * 1.5 for i in range(len(num_list1))]


rects1 = plt.bar(left=x, height=num_list1, width=0.4, alpha=0.8, color='salmon', label="Random")
rects2 = plt.bar(left=[i + 0.4 for i in x], height=num_list2, width=0.4, color='wheat', label="BDG")
rects3 = plt.bar(left=[i + 0.8 for i in x], height=num_list3, width=0.4, color='lightblue', label="SHP")
plt.ylim(0, 6)
plt.ylabel("Average Fanout")


plt.xticks([index + 0.6 for index in x], label_list)
plt.xlabel("Bucket Number")
plt.title("Fanout of running on different bucket numbers")
plt.legend()

for rect in rects1:
    height = rect.get_height()
    plt.text(rect.get_x() + rect.get_width() / 2, height+0.1, str(height), ha="center", va="bottom")
for rect in rects2:
    height = rect.get_height()
    plt.text(rect.get_x() + rect.get_width() / 2, height+0.1, str(height), ha="center", va="bottom")
for rect in rects3:
    height = rect.get_height()
    plt.text(rect.get_x() + rect.get_width() / 2, height+0.1, str(height), ha="center", va="bottom")
plt.savefig("barGraph_lj.png")
plt.show()
# import matplotlib.pyplot as plt
# import matplotlib
# matplotlib.rcParams['font.sans-serif'] = ['SimHei']
# matplotlib.rcParams['axes.unicode_minus'] = False

# label_list = ['2', '5', '10']
# num_list1 = [1.72, 3.31, 5.11]
# num_list2 = [1.47, 2.78, 3.62]
# num_list3 = [1.38, 1.95, 1.49]
# x = range(len(num_list1))
# x = [i * 1.5 for i in range(len(num_list1))]


# rects1 = plt.bar(left=x, height=num_list1, width=0.4, alpha=0.8, color='salmon', label="Random")
# rects2 = plt.bar(left=[i + 0.4 for i in x], height=num_list2, width=0.4, color='wheat', label="BDG")
# rects3 = plt.bar(left=[i + 0.8 for i in x], height=num_list3, width=0.4, color='lightblue', label="SHP")
# plt.ylim(0, 4)
# plt.ylabel("Average Fanout")


# plt.xticks([index + 0.6 for index in x], label_list)
# plt.xlabel("Bucket Number")
# plt.title("Fanout of running on different bucket numbers")
# plt.legend()

# for rect in rects1:
#     height = rect.get_height()
#     plt.text(rect.get_x() + rect.get_width() / 2, height+0.1, str(height), ha="center", va="bottom")
# for rect in rects2:
#     height = rect.get_height()
#     plt.text(rect.get_x() + rect.get_width() / 2, height+0.1, str(height), ha="center", va="bottom")
# for rect in rects3:
#     height = rect.get_height()
#     plt.text(rect.get_x() + rect.get_width() / 2, height+0.1, str(height), ha="center", va="bottom")
# plt.savefig("barGraph_youtube.png")
# plt.show()
