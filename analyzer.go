package inspector

import (
	"math"
	"encoding/json"
	"fmt"
	"log"
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

	everything := Harvest(repoName)


	if everything == nil {
		return returnBlocks
	}


	log.Printf("%v", everything)
	//	absoluteTotalNumberOfCommits := len(everything.Commits)

	totalNumberOfCommits := len(everything.Commits)
	for _, file := range everything.Files {


		for _, unit := range file.Units {

			numberOfCommits := unit.TimesChanged

//			//BULLSHIT START
//			rand.Seed(time.Now().UTC().UnixNano())
//
//			numberOfCommits = rand.Intn(100) + 1
//			totalNumberOfCommits = rand.Intn(30) + numberOfCommits
//			//BULLSHIT END
//
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
