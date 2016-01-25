package inspector

import (
	"math"
	"encoding/json"
	"fmt"
	"math/rand"
	"log"
	"time"
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

		for _, unit := range file.Units {

			numberOfCommits := len(unit.Commits)

			//BULLSHIT START
			rand.Seed(time.Now().UTC().UnixNano())

			numberOfCommits = rand.Intn(100) + 1
			totalNumberOfCommits = rand.Intn(30) + numberOfCommits
			//BULLSHIT END

			log.Println("-----")
			log.Printf("%d", totalNumberOfCommits)
			log.Printf("%d", numberOfCommits)
			log.Printf("%f", unit.RatioSum / float64(numberOfCommits))
			log.Printf("%f", float64(totalNumberOfCommits) / float64(numberOfCommits))
			log.Printf("%f", math.Log(float64(totalNumberOfCommits / numberOfCommits)))

			returnBlocks = append(returnBlocks, block{
				key: unit.Name,
				file: file.Path,
				unit_type: unit.Type,
				value: (unit.RatioSum / float64(numberOfCommits)) * math.Log(float64(totalNumberOfCommits) / float64(numberOfCommits)),
				size: unit.Size(),
			})
		}
	}

	return returnBlocks
}