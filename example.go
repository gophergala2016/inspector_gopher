package inspector
import "time"

func example() *Everything {
	var example map[string]*File = make(map[string]*File)
	example["onefile.go"] = &File{
		Path: "onefile.go",
		NumberOfLines: 300,
		Units: []*Unit{
			&Unit{
				Type: 0,
				Name: "makefunc",

				LineStart: 12,
				LineEnd: 45,

				RatioSum: 4.5,
				TimesChanged: 7,

				Commits: []*Commit{},
				File: nil,
			},
			&Unit{
				Type: 0,
				Name: "makefunc2",

				LineStart: 13,
				LineEnd: 23,

				RatioSum: 4.5,
				TimesChanged: 7,

				Commits: []*Commit{},
				File: nil,
			},
		},
	}

	commits := map[string]*Commit {
		"hash": &Commit{
			Contributor: nil, //Milenko
			Hash: "hash",
			Message: "Msg",
			Time: time.Now(),
			Units: []*Unit{
				//Uniti, ista struktura ako treba
			},
		},
	}

	return &Everything{
		Files: example,
		Commits: commits,
	}
}