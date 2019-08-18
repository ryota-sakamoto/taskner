package main

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type JobConfig struct {
	BeforeScript []Command
	AfterScript  []Command
	WatchFiles   []string
	Jobs         []Job
}

func read(file_name string) (*JobConfig, error) {
	file, err := os.Open(file_name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	m := map[string]interface{}{}
	err = yaml.Unmarshal([]byte(d), &m)
	if err != nil {
		return nil, err
	}

	before := []Command{}
	after := []Command{}
	watch_files := []string{}
	jobs := []Job{}

	if l, ok := m["before"]; ok {
		re, ok := l.([]interface{})
		if ok {
			for _, v := range re {
				v := v.(string)
				before = append(before, getCommand(v))
			}
		}
	}

	if l, ok := m["after"]; ok {
		re, ok := l.([]interface{})
		if ok {
			for _, v := range re {
				v := v.(string)
				after = append(after, getCommand(v))
			}
		}
	}

	if l, ok := m["watch"]; ok {
		re, ok := l.([]interface{})
		if ok {
			for _, v := range re {
				v := v.(string)
				watch_files = append(watch_files, v)
			}
		}
	}

	if l, ok := m["jobs"]; ok {
		re, ok := l.([]interface{})
		if ok {
			for _, v := range re {
				v := v.(map[interface{}]interface{})
				jobs = append(jobs, Job{
					Command:   getCommand(v["command"].(string)),
					StdinFile: v["stdin_from_file"].(string),
				})
			}
		}
	}

	j := JobConfig{
		BeforeScript: before,
		AfterScript:  after,
		WatchFiles:   watch_files,
		Jobs:         jobs,
	}

	return &j, nil
}
