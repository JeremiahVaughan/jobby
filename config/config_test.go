package config

import(
    "testing"
)

func Test_hasDuplicates(t *testing.T) {

    t.Run("no duplicates", func(t *testing.T) {
        input := []string{
            "a",
            "b",
        }
        got := hasDuplicates(input)
        if got {
            t.Errorf("error, did not expected duplicates but got")
        }
    })

    t.Run("has duplicates", func(t *testing.T) {
        input := []string{
            "b",
            "b",
        }
        got := hasDuplicates(input)
        if !got {
            t.Errorf("error, expected duplicates but did not get")
        }
    })


}
