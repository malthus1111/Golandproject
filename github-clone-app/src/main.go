package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"time"
	"encoding/csv"
	"log"
	"os/exec"
)

type Repository struct {
	Name        string    `json:"name"`
	UpdatedAt   time.Time `json:"updated_at"`
	CloneURL    string    `json:"clone_url"`
	Private     bool      `json:"private"`
	HTMLURL     string    `json:"html_url"`
	Description string    `json:"description"`
}

// Fonction pour écrire les informations des dépôts dans un fichier CSV.
func writeCSV(repos []Repository) {
	file, err := os.Create("./output/repositories.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Pour écrire l'en-tête du CSV.
	err = writer.Write([]string{"Name", "Updated At", "Clone URL", "Private", "URL", "Description"})
	if err != nil {
		log.Fatal(err)
	}

	// Pour écrire les données du dépôt dans le CSV.
	for _, repo := range repos {
		err := writer.Write([]string{repo.Name, repo.UpdatedAt.String(), repo.CloneURL, fmt.Sprintf("%v", repo.Private), repo.HTMLURL, repo.Description})
		if err != nil {
			log.Fatal(err)
		}
	}
}

func cloneRepo(cloneURL string, destination string) {
	cmd := exec.Command("git", "clone", cloneURL, destination)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
}

func gitPullLatestBranch(repoPath string) {
    
    err := os.Chdir(repoPath)
    if err != nil {
        log.Fatal(err)
    }

    
    cmd := exec.Command("git", "pull", "origin")
    err = cmd.Run()
    if err != nil {
        log.Fatal(err)
    }
}

func gitFetch(repoPath string) {
    // Naviguez vers le répertoire du dépôt cloné
    err := os.Chdir(repoPath)
    if err != nil {
        log.Fatal(err)
    }

    // Exécutez 'git fetch' pour récupérer toutes les références de branches
    cmd := exec.Command("git", "fetch", "--all")
    err = cmd.Run()
    if err != nil {
        log.Fatal(err)
    }
}


// ByUpdatedAt définit le type pour trier les dépôts par date de mise à jour.
type ByUpdatedAt []Repository

func (a ByUpdatedAt) Len() int           { return len(a) }
func (a ByUpdatedAt) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByUpdatedAt) Less(i, j int) bool { return a[i].UpdatedAt.After(a[j].UpdatedAt) }

func main() {
	fmt.Println("Starting the application...")
	token := os.Getenv("GITHUB_TOKEN")
	usernameOrOrgname := os.Getenv("USERNAME_OR_ORGNAME")

	// URL de l'API GitHub pour récupérer les dépôts d'un utilisateur
	page := 1
    perPage := 100
	url := fmt.Sprintf("https://api.github.com/users/%s/repos?page=%d&per_page=%d", usernameOrOrgname, page, perPage)

  for{
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Authorization", "token "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var repos []Repository
	json.Unmarshal(body, &repos)

	sort.Sort(ByUpdatedAt(repos))

	// Affichez les informations des dépôts.
	for _, repo := range repos {
		fmt.Printf("Name: %s, Updated At: %s, Clone URL: %s, Private: %v, URL: %s, Description: %s\n", repo.Name, repo.UpdatedAt, repo.CloneURL, repo.Private, repo.HTMLURL, repo.Description)
		// Cloner le dépôt
		fmt.Printf("Cloning repository %s\n", repo.Name)
	  cloneRepo(repo.CloneURL, "./output/"+repo.Name)
	  //pull
	  fmt.Printf("Pulling repository %s\n", repo.Name)
	  // gitPullLatestBranch("./output/" + repo.Name)
	  fmt.Printf("Done pulling repository %s\n", repo.Name)
	  //fetch
	  fmt.Printf("Fetching repository %s\n", repo.Name)
	  gitFetch("./output/" + repo.Name)
	  fmt.Printf("Done fetching repository %s\n", repo.Name)
	}

    // Pour écrire les informations des dépôts dans un fichier CSV.
	writeCSV(repos)

	if len(repos) < perPage {
		break
	}
	page++
	url = fmt.Sprintf("https://api.github.com/users/%s/repos?page=%d&per_page=%d", usernameOrOrgname, page, perPage)

  }
}
