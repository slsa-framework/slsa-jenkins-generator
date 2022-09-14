package provenance

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type DigestSet map[string]string

type Statement struct {
	Type          string    `json:"_type"`
	Subject       []Subject `json:"subject"`
	PredicateType string    `json:"predicateType"`
	Predicate     `json:"predicate"`
}

type Subject struct {
	Name   string    `json:"name"`
	Digest DigestSet `json:"digest"`
}

type Invocation struct {
	ConfigSource `json:"configSource"`
	Environment  JenkinsContext `json:"environment"`
}

type ConfigSource struct {
	URI        string    `json:"uri"`
	Digest     DigestSet `json:"digest"`
	EntryPoint string    `json:"entrypoint"`
}

type Predicate struct {
	BuildType  string `json:"buildtype"`
	Builder    `json:"builder"`
	Invocation `json:"invocation"`
	Materials  []Meterial `json:"materials"`
}

type Meterial struct {
	URI    string    `json:"uri"`
	Digest DigestSet `json:"digest"`
}

type Builder struct {
	Id string `json:"id"`
}

type JenkinsContext struct {
	Actor         string                 `json:"actor,omitempty"`
	Payload       map[string]interface{} `json:"event_payload,omitempty"`
	BuildURL      string                 `json:"build_url,omitempty"`
	RepositoryURL string                 `json:"payload_repository_git_url,omitempty"`
	SHA           string                 `json:"payload_after,omitempty"`
	EventName     string                 `json:"eventname,omitempty"`
}

func getSHA256(target string) string {

	f, err := os.Open(target)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}
	return hex.EncodeToString(h.Sum(nil))

}

func GenerateSLSA(artifact string, outpath string) {

	repo := os.Getenv("payload_repository_git_url") + "@" + os.Getenv("payload_ref")
	digest := DigestSet{"sha1": os.Getenv("payload_after")}

	//Generate Provenance
	prov := Statement{
		//Type Field
		Type: "https://in-toto.io/Statement/v0.1",
		//PredicateType Field
		PredicateType: "https://slsa.dev/provenance/v0.2",
	}

	//Subject Field
	var subjects []Subject
	name := strings.Split(artifact, "/")
	subjects = append(subjects, Subject{
		Name:   name[len(name)-1],
		Digest: DigestSet{"sha256": getSHA256(artifact)},
	})
	prov.Subject = subjects

	//Predicate Field
	var meterials []Meterial
	meterials = append(meterials, Meterial{
		URI:    repo,
		Digest: digest,
	})
	invoca := Invocation{}

	//Environment Field of Invocation
	var out bytes.Buffer
	ctx := JenkinsContext{}
	err := json.Indent(&out, []byte(os.Getenv("payload")), "", "\t")
	if err != nil {
		log.Fatal(err)
		return
	}

	err = json.Unmarshal([]byte("{\"event_payload\":"+out.String()+"}"), &ctx)
	if err != nil {
		log.Fatal(err)
		return
	}

	ctx.EventName = os.Getenv("x_github_event")
	ctx.Actor = os.Getenv("payload_head_commit_author_name")
	ctx.BuildURL = os.Getenv("BUILD_URL")
	invoca.Environment = ctx

	//ConfigSource Field of Invocation
	invoca.ConfigSource = ConfigSource{
		URI:        repo,
		Digest:     digest,
		EntryPoint: "config.xml",
	}

	prov.Predicate = Predicate{
		BuildType:  "https://www.jenkins.io/Pipeline", //Temporay, Not yet fixed
		Builder:    Builder{Id: ctx.BuildURL},
		Invocation: invoca,
		Materials:  meterials,
	}

	//Write to json file
	data, err := json.MarshalIndent(prov, "", "\t")
	if err != nil {
		log.Fatal(err)
		return
	}

	err = ioutil.WriteFile(outpath+"/provenance.slsa", data, 0644)
	if err != nil {
		log.Fatal(err)
	}

}
