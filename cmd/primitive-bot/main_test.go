package main

import (
	"os"
	"reflect"
	"testing"

	"github.com/lazy-void/primitive-bot/pkg/primitive"
	"github.com/lazy-void/primitive-bot/pkg/queue"
)

func TestRestoreQueue(t *testing.T) {
	workers := 1
	tests := []struct {
		name       string
		logData    string
		operations []queue.Operation
	}{
		{
			name: "One operation in the queue",
			logData: `
INFO	2021/05/23 16:00:42 Starting to listen for updates...
INFO	2021/05/23 16:00:48 Message: text '' from the user 'Kir' with the ID '295434263'
INFO	2021/05/23 16:00:49 Callback Query: data '/create' from the user 'Kir' with the ID '295434263'
INFO	2021/05/23 16:00:49 Enqueued: user id 295434263 | input inputs/AQADntiNoi4AAwSIAgAB.jpg | iterations=200, shape=0, alpha=128, repeat=1, resolution=1280, extension=jpg
INFO	2021/05/23 16:00:49 Creating: user id 295434263 | input inputs/AQADntiNoi4AAwSIAgAB.jpg | output outputs/295434263_1621774849.jpg | iterations=200, shape=0, alpha=128, repeat=1, resolution=1280, extension=jpg
`,
			operations: []queue.Operation{
				{
					UserID:  295434263,
					ImgPath: "inputs/AQADntiNoi4AAwSIAgAB.jpg",
					Config:  primitive.New(workers),
				},
			},
		},
		{
			name: "One operaion that was finished",
			logData: `
INFO	2021/05/23 16:00:42 Starting to listen for updates...
INFO	2021/05/23 16:00:48 Message: text '' from the user 'Kir' with the ID '295434263'
INFO	2021/05/23 16:00:49 Callback Query: data '/create' from the user 'Kir' with the ID '295434263'
INFO	2021/05/23 16:00:49 Enqueued: user id 295434263 | input inputs/AQADntiNoi4AAwSIAgAB.jpg | iterations=200, shape=0, alpha=128, repeat=1, resolution=1280, extension=jpg
INFO	2021/05/23 16:00:49 Creating: user id 295434263 | input inputs/AQADntiNoi4AAwSIAgAB.jpg | output outputs/295434263_1621774849.jpg | iterations=200, shape=0, alpha=128, repeat=1, resolution=1280, extension=jpg
INFO	2021/05/23 16:00:52 Finished: user id 295434263 | input inputs/AQADntiNoi4AAwSIAgAB.jpg | output outputs/295434263_1621774849.jpg | 3.1 seconds
INFO	2021/05/23 16:00:55 Sent: user id 295434263 | output outputs/295434263_1621774849.jpg
`,
			operations: []queue.Operation{},
		},
		{
			name: "Three operations, one was finished",
			logData: `
INFO	2021/05/23 16:45:06 Starting to listen for updates...
INFO	2021/05/23 16:45:21 Message: text '' from the user 'Kir' with the ID '295434263'
INFO	2021/05/23 16:45:23 Callback Query: data '/create' from the user 'Kir' with the ID '295434263'
INFO	2021/05/23 16:45:23 Enqueued: user id 295434263 | input inputs/AQADntiNoi4AAwSIAgAB.jpg | iterations=10, shape=0, alpha=128, repeat=1, resolution=1280, extension=jpg
INFO	2021/05/23 16:45:23 Creating: user id 295434263 | input inputs/AQADntiNoi4AAwSIAgAB.jpg | output outputs/295434263_1621777523.jpg | iterations=200, shape=0, alpha=128, repeat=1, resolution=1280, extension=jpg
INFO	2021/05/23 16:45:24 Callback Query: data '/create' from the user 'Kir' with the ID '295434263'
INFO	2021/05/23 16:45:24 Enqueued: user id 295434263 | input inputs/AQADntiNoi4AAwSIAgAB.jpg | iterations=20, shape=2, alpha=128, repeat=1, resolution=1280, extension=jpg
INFO	2021/05/23 16:45:25 Callback Query: data '/create' from the user 'Kir' with the ID '2954342632
INFO	2021/05/23 16:45:25 Enqueued: user id 295434263 | input inputs/AQADntiNoi4AAwSIAgAB.jpg | iterations=30, shape=1, alpha=128, repeat=1, resolution=1280, extension=jpg
INFO	2021/05/23 16:49:12 Finished: user id 295434263 | input inputs/AQADntiNoi4AAwSIAgAB.jpg | output outputs/295434263_1621774849.jpg | 5.1 seconds
INFO	2021/05/23 16:49:20 Sent: user id 295434263 | output outputs/295434263_1621774849.jpg
`,
			operations: []queue.Operation{
				{
					UserID:  295434263,
					ImgPath: "inputs/AQADntiNoi4AAwSIAgAB.jpg",
					Config: func() primitive.Config {
						c := primitive.New(workers)
						c.Iterations = 20
						c.Shape = 2
						return c
					}(),
				},
				{
					UserID:  295434263,
					ImgPath: "inputs/AQADntiNoi4AAwSIAgAB.jpg",
					Config: func() primitive.Config {
						c := primitive.New(workers)
						c.Iterations = 30
						c.Shape = 1
						return c
					}(),
				},
			},
		},
		{
			name: "No operations in the queue",
			logData: `
INFO	2021/05/23 16:00:42 Starting to listen for updates...
INFO	2021/05/23 16:00:48 Message: text '' from the user 'Kir' with the ID '295434263'
`,
			operations: []queue.Operation{},
		},
		{
			name: "All operations are finished",
			logData: `
INFO	2021/05/23 16:54:10 Starting to listen for updates...
INFO	2021/05/23 16:54:10 Message: text '' from the user 'Kir' with the ID '295434263'
INFO	2021/05/23 16:54:14 Callback Query: data '/iter' from the user 'Kir' with the ID '295434263'
INFO	2021/05/23 16:54:16 Callback Query: data '/iter/input' from the user 'Kir' with the ID '295434263'
INFO	2021/05/23 16:54:18 Message: text '10' from the user 'Kir' with the ID '295434263'
INFO	2021/05/23 16:54:20 Callback Query: data '/' from the user 'Kir' with the ID '295434263'
INFO	2021/05/23 16:54:21 Callback Query: data '/create' from the user 'Kir' with the ID '295434263'
INFO	2021/05/23 16:54:21 Enqueued: user id 295434263 | input inputs/AQADntiNoi4AAwSIAgAB.jpg | iterations=10, shape=0, alpha=128, repeat=1, resolution=1280, extension=jpg
INFO	2021/05/23 16:54:22 Creating: user id 295434263 | input inputs/AQADntiNoi4AAwSIAgAB.jpg | output outputs/295434263_1621778062.jpg | iterations=10, shape=0, alpha=128, repeat=1, resolution=1280, extension=jpg
INFO	2021/05/23 16:54:22 Callback Query: data '/shape' from the user 'Kir' with the ID '295434263'
INFO	2021/05/23 16:54:24 Finished: user id 295434263 | input inputs/AQADntiNoi4AAwSIAgAB.jpg | output outputs/295434263_1621778062.jpg | 1.7 seconds
INFO	2021/05/23 16:54:24 Sent: user id 295434263 | output outputs/295434263_1621778062.jpg
INFO	2021/05/23 16:54:25 Callback Query: data '/shape/1' from the user 'Kir' with the ID '295434263'
INFO	2021/05/23 16:54:26 Callback Query: data '/' from the user 'Kir' with the ID '295434263'
INFO	2021/05/23 16:54:27 Callback Query: data '/create' from the user 'Kir' with the ID '295434263'
INFO	2021/05/23 16:54:27 Enqueued: user id 295434263 | input inputs/AQADntiNoi4AAwSIAgAB.jpg | iterations=10, shape=1, alpha=128, repeat=1, resolution=1280, extension=jpg
INFO	2021/05/23 16:54:28 Creating: user id 295434263 | input inputs/AQADntiNoi4AAwSIAgAB.jpg | output outputs/295434263_1621778068.jpg | iterations=10, shape=1, alpha=128, repeat=1, resolution=1280, extension=jpg
INFO	2021/05/23 16:54:28 Callback Query: data '/shape' from the user 'Kir' with the ID '295434263'
INFO	2021/05/23 16:54:29 Finished: user id 295434263 | input inputs/AQADntiNoi4AAwSIAgAB.jpg | output outputs/295434263_1621778068.jpg | 1.4 seconds
INFO	2021/05/23 16:54:30 Callback Query: data '/shape/6' from the user 'Kir' with the ID '295434263'
INFO	2021/05/23 16:54:30 Sent: user id 295434263 | output outputs/295434263_1621778068.jpg
INFO	2021/05/23 16:54:31 Callback Query: data '/' from the user 'Kir' with the ID '295434263'
INFO	2021/05/23 16:54:32 Callback Query: data '/create' from the user 'Kir' with the ID '295434263'
INFO	2021/05/23 16:54:32 Enqueued: user id 295434263 | input inputs/AQADntiNoi4AAwSIAgAB.jpg | iterations=10, shape=6, alpha=128, repeat=1, resolution=1280, extension=jpg
INFO	2021/05/23 16:54:32 Creating: user id 295434263 | input inputs/AQADntiNoi4AAwSIAgAB.jpg | output outputs/295434263_1621778072.jpg | iterations=10, shape=6, alpha=128, repeat=1, resolution=1280, extension=jpg
INFO	2021/05/23 16:54:33 Finished: user id 295434263 | input inputs/AQADntiNoi4AAwSIAgAB.jpg | output outputs/295434263_1621778072.jpg | 1.0 seconds
INFO	2021/05/23 16:54:33 Sent: user id 295434263 | output outputs/295434263_1621778072.jpg
`,
			operations: []queue.Operation{},
		},
	}

	logPath := "test_log.txt"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create fake log
			err := os.WriteFile(logPath, []byte(tt.logData), 0600)
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(logPath)

			// create queue
			q := queue.New()
			expected := queue.New()
			for _, op := range tt.operations {
				expected.Enqueue(op)
			}

			// restore queue
			if err = restoreQueue(logPath, q, workers); err != nil {
				t.Error(err)
			}

			// check
			if !reflect.DeepEqual(expected, q) {
				t.Error("queues are different")
			}
		})
	}
}
