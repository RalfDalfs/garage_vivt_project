package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const banner = "Команды: + НОМЕР, - НОМЕР, list, exit\n"

func TestCLI(t *testing.T) {
	cases := []struct{ name, input, want string }{
		{"empty", "", ""},
		{"blank", "  \n\t\nexit\n", ""},
		{"empty_list", "list\n", "Стоянка от дальнего конца к выезду:\n"},
		{"one_car", "+ А101\n- А101\n", "Прибыл А101\nУехал А101: перемещений 0\n"},
		{"unknown", "- X\n", "Ошибка: автомобиль не найден\n"},
		{"duplicate", "+ A\n+ A\nlist\n", "Прибыл A\nОшибка: автомобиль уже на стоянке\nСтоянка от дальнего конца к выезду:\nA: перемещений 0\n"},
		{"exit", "exit\n+ A\n", ""},
		{"eof_without_newline", "+ A", "Прибыл A\n"},
		{"invalid", "+\n- A B\nlist A\nexit A\n?\n", strings.Repeat("Ошибка: используйте + НОМЕР, - НОМЕР, list или exit\n", 5)},
		{"nearest", "+ A\n+ B\n- B\nlist\n", "Прибыл A\nПрибыл B\nУехал B: перемещений 0\nСтоянка от дальнего конца к выезду:\nA: перемещений 0\n"},
		{"reset_on_new_visit", "+ A\n+ B\n- A\n- B\n+ B\n- B\n", "Прибыл A\nПрибыл B\nВременно выехал B\nУехал A: перемещений 0\nВернулся B\nУехал B: перемещений 1\nПрибыл B\nУехал B: перемещений 0\n"},
	}
	input, err := os.ReadFile("../testdata/demo.txt")
	if err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile("../testdata/demo_expected.txt")
	if err != nil {
		t.Fatal(err)
	}
	cases = append(cases, struct{ name, input, want string }{"order_and_accumulation", string(input), strings.TrimPrefix(string(want), banner)})
	for _, version := range []string{"linear", "modular"} {
		t.Run(version, func(t *testing.T) {
			binary := filepath.Join(t.TempDir(), version)
			if out, err := exec.Command("go", "build", "-o", binary, "../cmd/"+version).CombinedOutput(); err != nil {
				t.Fatalf("build: %s %v", out, err)
			}
			for _, tc := range cases {
				t.Run(tc.name, func(t *testing.T) {
					cmd := exec.Command(binary)
					cmd.Stdin = strings.NewReader(tc.input)
					out, err := cmd.CombinedOutput()
					if err != nil {
						t.Fatal(err)
					}
					if string(out) != banner+tc.want {
						t.Fatalf("got:\n%s\nwant:\n%s", out, banner+tc.want)
					}
				})
			}
			t.Run("too_long_line", func(t *testing.T) {
				cmd := exec.Command(binary)
				cmd.Stdin = strings.NewReader(strings.Repeat("X", 70000))
				out, err := cmd.CombinedOutput()
				if err == nil || !strings.Contains(string(out), "Ошибка чтения:") {
					t.Fatalf("%s %v", out, err)
				}
			})
		})
	}
}
