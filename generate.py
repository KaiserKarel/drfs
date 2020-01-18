
def generate_sequence(filename, lines, segments):
    with open(filename, "w+") as f:
        for _ in range(lines):
            for s in segments:
                f.write(s[0]*s[1])

if __name__ == "__main__":
    generate_sequence("testdata/alr.gen.txt", 50, [("a", 4092), ("l", 2), ("r", 2), ])