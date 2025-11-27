window.BENCHMARK_DATA = {
  "lastUpdate": 1764252679072,
  "repoUrl": "https://github.com/y-l-g/pogo",
  "entries": {
    "Benchmark": [
      {
        "commit": {
          "author": {
            "email": "youenn.legouedec@gmail.com",
            "name": "y-l-g",
            "username": "y-l-g"
          },
          "committer": {
            "email": "youenn.legouedec@gmail.com",
            "name": "y-l-g",
            "username": "y-l-g"
          },
          "distinct": true,
          "id": "9755a49ac028922950854c8985e907cbf6af9f45",
          "message": "fix readme title",
          "timestamp": "2025-11-27T14:51:40+01:00",
          "tree_id": "2774dc86d19c198da4de8bd363d99b6471e5f6bc",
          "url": "https://github.com/y-l-g/pogo/commit/9755a49ac028922950854c8985e907cbf6af9f45"
        },
        "date": 1764251556604,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkAllocate",
            "value": 91.3,
            "unit": "ns/op\t      24 B/op\t       1 allocs/op",
            "extra": "11572376 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - ns/op",
            "value": 91.3,
            "unit": "ns/op",
            "extra": "11572376 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - B/op",
            "value": 24,
            "unit": "B/op",
            "extra": "11572376 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "11572376 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel",
            "value": 126.9,
            "unit": "ns/op\t      24 B/op\t       1 allocs/op",
            "extra": "9394041 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - ns/op",
            "value": 126.9,
            "unit": "ns/op",
            "extra": "9394041 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - B/op",
            "value": 24,
            "unit": "B/op",
            "extra": "9394041 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "9394041 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt",
            "value": 70.22,
            "unit": "ns/op\t58334.77 MB/s\t       0 B/op\t       0 allocs/op",
            "extra": "17063862 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - ns/op",
            "value": 70.22,
            "unit": "ns/op",
            "extra": "17063862 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - MB/s",
            "value": 58334.77,
            "unit": "MB/s",
            "extra": "17063862 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "17063862 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "17063862 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON",
            "value": 3569,
            "unit": "ns/op\t    1360 B/op\t      36 allocs/op",
            "extra": "333613 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - ns/op",
            "value": 3569,
            "unit": "ns/op",
            "extra": "333613 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - B/op",
            "value": 1360,
            "unit": "B/op",
            "extra": "333613 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - allocs/op",
            "value": 36,
            "unit": "allocs/op",
            "extra": "333613 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack",
            "value": 680.1,
            "unit": "ns/op\t     192 B/op\t       1 allocs/op",
            "extra": "1762215 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - ns/op",
            "value": 680.1,
            "unit": "ns/op",
            "extra": "1762215 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - B/op",
            "value": 192,
            "unit": "B/op",
            "extra": "1762215 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1762215 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation",
            "value": 433.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "2744976 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - ns/op",
            "value": 433.4,
            "unit": "ns/op",
            "extra": "2744976 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "2744976 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "2744976 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus",
            "value": 731.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1691247 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - ns/op",
            "value": 731.6,
            "unit": "ns/op",
            "extra": "1691247 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1691247 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1691247 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "youenn.legouedec@gmail.com",
            "name": "y-l-g",
            "username": "y-l-g"
          },
          "committer": {
            "email": "youenn.legouedec@gmail.com",
            "name": "y-l-g",
            "username": "y-l-g"
          },
          "distinct": true,
          "id": "b901183ef3992b3388372c9a39956dedfc200460",
          "message": "wip",
          "timestamp": "2025-11-27T15:10:18+01:00",
          "tree_id": "28c25deb6b5f4fcf4dd87bd813b70b6ff9d5179d",
          "url": "https://github.com/y-l-g/pogo/commit/b901183ef3992b3388372c9a39956dedfc200460"
        },
        "date": 1764252678180,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkAllocate",
            "value": 90.22,
            "unit": "ns/op\t      24 B/op\t       1 allocs/op",
            "extra": "11685020 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - ns/op",
            "value": 90.22,
            "unit": "ns/op",
            "extra": "11685020 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - B/op",
            "value": 24,
            "unit": "B/op",
            "extra": "11685020 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "11685020 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel",
            "value": 111.7,
            "unit": "ns/op\t      24 B/op\t       1 allocs/op",
            "extra": "10596259 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - ns/op",
            "value": 111.7,
            "unit": "ns/op",
            "extra": "10596259 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - B/op",
            "value": 24,
            "unit": "B/op",
            "extra": "10596259 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "10596259 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt",
            "value": 54.68,
            "unit": "ns/op\t74908.63 MB/s\t       0 B/op\t       0 allocs/op",
            "extra": "21814486 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - ns/op",
            "value": 54.68,
            "unit": "ns/op",
            "extra": "21814486 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - MB/s",
            "value": 74908.63,
            "unit": "MB/s",
            "extra": "21814486 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "21814486 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "21814486 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON",
            "value": 3567,
            "unit": "ns/op\t    1360 B/op\t      36 allocs/op",
            "extra": "327739 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - ns/op",
            "value": 3567,
            "unit": "ns/op",
            "extra": "327739 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - B/op",
            "value": 1360,
            "unit": "B/op",
            "extra": "327739 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - allocs/op",
            "value": 36,
            "unit": "allocs/op",
            "extra": "327739 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack",
            "value": 737.5,
            "unit": "ns/op\t     192 B/op\t       1 allocs/op",
            "extra": "1697911 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - ns/op",
            "value": 737.5,
            "unit": "ns/op",
            "extra": "1697911 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - B/op",
            "value": 192,
            "unit": "B/op",
            "extra": "1697911 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1697911 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation",
            "value": 408.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "2919946 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - ns/op",
            "value": 408.8,
            "unit": "ns/op",
            "extra": "2919946 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "2919946 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "2919946 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus",
            "value": 682.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1691425 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - ns/op",
            "value": 682.2,
            "unit": "ns/op",
            "extra": "1691425 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1691425 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1691425 times\n4 procs"
          }
        ]
      }
    ]
  }
}