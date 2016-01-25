package inspector

import (
	"math"
	"encoding/json"
	"fmt"
)

type block struct {
	key       string
	file      string
	unit_type int
	value     float64
	size      int
}

func (b block) MarshalJSON() ([]byte, error) {
	m := make(map[string]string)

	m["key"] = b.key
	m["file"] = b.file
	m["unit_type"] = fmt.Sprintf("%d", b.unit_type)
	m["value"] = fmt.Sprintf("%.6f", b.value)
	m["size"] = fmt.Sprintf("%d", b.size)

	return json.Marshal(m)
}

func AnalyzeRepo(repoName string) (returnBlocks []block) {
	//	everything := Harvest(repoName)

	everything := example()


	//	absoluteTotalNumberOfCommits := len(everything.Commits)

	for _, file := range everything.Files {

		totalNumberOfCommits := len(file.Commits)

		if totalNumberOfCommits == 0 { totalNumberOfCommits = 1}

		for _, unit := range file.Units {

			numberOfCommits := len(unit.Commits)

			if numberOfCommits == 0 { numberOfCommits = 1}

			returnBlocks = append(returnBlocks, block{
				key: unit.Name,
				file: file.Path,
				unit_type: unit.Type,
				value: (unit.RatioSum / float64(numberOfCommits)) * math.Log(float64(totalNumberOfCommits / numberOfCommits)),
				size: unit.Size(),
			})
		}
	}

	return returnBlocks
}
