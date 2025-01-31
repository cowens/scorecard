// Copyright 2021 Security Scorecard Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package data

import (
	"embed"
	"encoding/csv"
	"errors"
	"fmt"
	"io"

	"github.com/jszwec/csvutil"

	"github.com/ossf/scorecard/repos"
)

//go:embed projects.csv
var csvFile embed.FS

type Iterator interface {
	HasNext() bool
	Next() (repos.RepoURL, error)
}

func MakeIterator() (Iterator, error) {
	csvReader, err := csvFile.Open("projects.csv")
	if err != nil {
		return nil, fmt.Errorf("error in Open: %w", err)
	}
	return MakeIteratorFrom(csvReader)
}

func MakeIteratorFrom(reader io.Reader) (Iterator, error) {
	dec, err := csvutil.NewDecoder(csv.NewReader(reader))
	if err != nil {
		return nil, fmt.Errorf("error in csvutil.NewDecoder: %w", err)
	}
	return &csvIterator{decoder: dec}, nil
}

type inputRepo struct {
	Repo     string `csv:"repo"`
	Metadata string `csv:"metadata"`
}

type csvIterator struct {
	decoder *csvutil.Decoder
	err     error
	next    inputRepo
}

func (reader *csvIterator) HasNext() bool {
	reader.err = reader.decoder.Decode(&reader.next)
	return !errors.Is(reader.err, io.EOF)
}

func (reader *csvIterator) Next() (repos.RepoURL, error) {
	if reader.err != nil {
		return repos.RepoURL{}, reader.err
	}
	ret := repos.RepoURL{}
	var err error
	err = ret.Set(reader.next.Repo)
	if err == nil {
		err = ret.ValidGitHubURL()
	}
	return ret, errors.Unwrap(err)
}
