package crack

import "testing"

func TestXxx(t *testing.T) {
	s1 := Steps{
		Xpath: `//*[@id="loginForm"]/fieldset/div[2]/input`,
		Dict:  []string{"druid", "admin"},
	}
	s2 := Steps{
		Xpath: `//*[@id="loginForm"]/fieldset/div[3]/input`,
		Dict:  []string{"druid", "admin"},
	}
	s3 := Steps{
		Xpath: `//*[@id="loginBtn"]`,
	}
	var steps []Steps
	steps = append(steps, s1, s2, s3)
	CheckLogin("", [3]Steps(steps))
}
