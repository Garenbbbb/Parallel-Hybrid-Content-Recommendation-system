import subprocess
import sys

def main():
  seq_command = "go run recommend.go random 1 1"
  times = 1
  base = 0
  for i in range(times):
    process = subprocess.Popen(seq_command, stdout=subprocess.PIPE, stderr=subprocess.PIPE, shell=True, universal_newlines=True)
    output_lines = process.stdout.readlines()
    base += float(output_lines[-1].strip()[:6])
    process.wait()
    if process.returncode != 0:
        print("Error: Go program returned a non-zero exit code.")
  base = base/times
  thread = [2, 4, 6, 8, 12]
  go_command = 'go run recommend.go random '

  data = []

  for k in thread:
    temp = 0
    for i in range(times):
      process = subprocess.Popen(go_command + str(k) + " " + str(k), stdout=subprocess.PIPE, stderr=subprocess.PIPE, shell=True, universal_newlines=True)
      output_lines = process.stdout.readlines()
      temp += float(output_lines[-1].strip()[:6])
      process.wait()
      if process.returncode != 0:
          print("Error: Go program returned a non-zero exit code.")
    data.append(round(base/temp/times, 3))
        

  import matplotlib.pyplot as plt


  print(data)

  # Create a line graph
  plt.plot(thread, data)
  plt.xlabel('threads')
  plt.ylabel('Speed-up')
  plt.title(" test_case")
  # plt.show()
  plt.savefig("analysis.png")

if __name__ == "__main__":
  main()