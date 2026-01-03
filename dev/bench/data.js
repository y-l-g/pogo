window.BENCHMARK_DATA = {
  "lastUpdate": 1767480921740,
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
          "id": "a76c192194ace862fbd297e67c7df5d8bc3cad4f",
          "message": "wip",
          "timestamp": "2025-11-27T15:50:35+01:00",
          "tree_id": "9d3d2290825252b8720139eaee0866b110da8bd1",
          "url": "https://github.com/y-l-g/pogo/commit/a76c192194ace862fbd297e67c7df5d8bc3cad4f"
        },
        "date": 1764255171843,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkAllocate",
            "value": 82.8,
            "unit": "ns/op\t      24 B/op\t       1 allocs/op",
            "extra": "15338876 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - ns/op",
            "value": 82.8,
            "unit": "ns/op",
            "extra": "15338876 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - B/op",
            "value": 24,
            "unit": "B/op",
            "extra": "15338876 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "15338876 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel",
            "value": 116,
            "unit": "ns/op\t      24 B/op\t       1 allocs/op",
            "extra": "10301402 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - ns/op",
            "value": 116,
            "unit": "ns/op",
            "extra": "10301402 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - B/op",
            "value": 24,
            "unit": "B/op",
            "extra": "10301402 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "10301402 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt",
            "value": 55.5,
            "unit": "ns/op\t73797.46 MB/s\t       0 B/op\t       0 allocs/op",
            "extra": "18158590 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - ns/op",
            "value": 55.5,
            "unit": "ns/op",
            "extra": "18158590 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - MB/s",
            "value": 73797.46,
            "unit": "MB/s",
            "extra": "18158590 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "18158590 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "18158590 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON",
            "value": 3510,
            "unit": "ns/op\t    1360 B/op\t      36 allocs/op",
            "extra": "335815 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - ns/op",
            "value": 3510,
            "unit": "ns/op",
            "extra": "335815 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - B/op",
            "value": 1360,
            "unit": "B/op",
            "extra": "335815 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - allocs/op",
            "value": 36,
            "unit": "allocs/op",
            "extra": "335815 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack",
            "value": 714.4,
            "unit": "ns/op\t     192 B/op\t       1 allocs/op",
            "extra": "1728404 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - ns/op",
            "value": 714.4,
            "unit": "ns/op",
            "extra": "1728404 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - B/op",
            "value": 192,
            "unit": "B/op",
            "extra": "1728404 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1728404 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation",
            "value": 403.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "2947850 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - ns/op",
            "value": 403.3,
            "unit": "ns/op",
            "extra": "2947850 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "2947850 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "2947850 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus",
            "value": 698.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1728141 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - ns/op",
            "value": 698.5,
            "unit": "ns/op",
            "extra": "1728141 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1728141 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1728141 times\n4 procs"
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
          "id": "6045ca61cddc8b706052209a1d7f15eb2a0f0388",
          "message": "wip",
          "timestamp": "2025-11-27T23:04:08+01:00",
          "tree_id": "dacef6a228fc20cd6e29371684eb8d5f3e12a8b1",
          "url": "https://github.com/y-l-g/pogo/commit/6045ca61cddc8b706052209a1d7f15eb2a0f0388"
        },
        "date": 1764281102192,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkAllocate",
            "value": 87.94,
            "unit": "ns/op\t      24 B/op\t       1 allocs/op",
            "extra": "12647845 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - ns/op",
            "value": 87.94,
            "unit": "ns/op",
            "extra": "12647845 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - B/op",
            "value": 24,
            "unit": "B/op",
            "extra": "12647845 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "12647845 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel",
            "value": 114.6,
            "unit": "ns/op\t      24 B/op\t       1 allocs/op",
            "extra": "10338998 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - ns/op",
            "value": 114.6,
            "unit": "ns/op",
            "extra": "10338998 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - B/op",
            "value": 24,
            "unit": "B/op",
            "extra": "10338998 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "10338998 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt",
            "value": 54.8,
            "unit": "ns/op\t74740.58 MB/s\t       0 B/op\t       0 allocs/op",
            "extra": "21927301 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - ns/op",
            "value": 54.8,
            "unit": "ns/op",
            "extra": "21927301 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - MB/s",
            "value": 74740.58,
            "unit": "MB/s",
            "extra": "21927301 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "21927301 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "21927301 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON",
            "value": 3495,
            "unit": "ns/op\t    1360 B/op\t      36 allocs/op",
            "extra": "338536 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - ns/op",
            "value": 3495,
            "unit": "ns/op",
            "extra": "338536 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - B/op",
            "value": 1360,
            "unit": "B/op",
            "extra": "338536 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - allocs/op",
            "value": 36,
            "unit": "allocs/op",
            "extra": "338536 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack",
            "value": 705.3,
            "unit": "ns/op\t     192 B/op\t       1 allocs/op",
            "extra": "1701756 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - ns/op",
            "value": 705.3,
            "unit": "ns/op",
            "extra": "1701756 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - B/op",
            "value": 192,
            "unit": "B/op",
            "extra": "1701756 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1701756 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation",
            "value": 405.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "2970650 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - ns/op",
            "value": 405.7,
            "unit": "ns/op",
            "extra": "2970650 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "2970650 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "2970650 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus",
            "value": 702.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1724515 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - ns/op",
            "value": 702.1,
            "unit": "ns/op",
            "extra": "1724515 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1724515 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1724515 times\n4 procs"
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
          "id": "71dd95396e5a23c3b35c2c53f50b8dad295139ae",
          "message": "wip",
          "timestamp": "2025-11-28T13:32:48+01:00",
          "tree_id": "1e0bd6850d065d65fa57182606d1a3f7b8b3e7be",
          "url": "https://github.com/y-l-g/pogo/commit/71dd95396e5a23c3b35c2c53f50b8dad295139ae"
        },
        "date": 1764333220329,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkAllocate",
            "value": 89.67,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "14196847 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - ns/op",
            "value": 89.67,
            "unit": "ns/op",
            "extra": "14196847 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "14196847 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "14196847 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel",
            "value": 113.5,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "10505529 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - ns/op",
            "value": 113.5,
            "unit": "ns/op",
            "extra": "10505529 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "10505529 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "10505529 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt",
            "value": 54.58,
            "unit": "ns/op\t75046.89 MB/s\t       0 B/op\t       0 allocs/op",
            "extra": "21929022 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - ns/op",
            "value": 54.58,
            "unit": "ns/op",
            "extra": "21929022 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - MB/s",
            "value": 75046.89,
            "unit": "MB/s",
            "extra": "21929022 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "21929022 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "21929022 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON",
            "value": 3586,
            "unit": "ns/op\t    1360 B/op\t      36 allocs/op",
            "extra": "330172 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - ns/op",
            "value": 3586,
            "unit": "ns/op",
            "extra": "330172 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - B/op",
            "value": 1360,
            "unit": "B/op",
            "extra": "330172 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - allocs/op",
            "value": 36,
            "unit": "allocs/op",
            "extra": "330172 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack",
            "value": 704.1,
            "unit": "ns/op\t     192 B/op\t       1 allocs/op",
            "extra": "1685114 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - ns/op",
            "value": 704.1,
            "unit": "ns/op",
            "extra": "1685114 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - B/op",
            "value": 192,
            "unit": "B/op",
            "extra": "1685114 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1685114 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation",
            "value": 404.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "3015831 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - ns/op",
            "value": 404.7,
            "unit": "ns/op",
            "extra": "3015831 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "3015831 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "3015831 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus",
            "value": 688.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1552785 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - ns/op",
            "value": 688.2,
            "unit": "ns/op",
            "extra": "1552785 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1552785 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1552785 times\n4 procs"
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
          "id": "87feb7e884dbef2ae5d6a42eae6cf04c482b49ae",
          "message": "wip",
          "timestamp": "2025-11-28T14:06:38+01:00",
          "tree_id": "0f31c1c61d7972a92c7b1e3a3b9c41218a8068b7",
          "url": "https://github.com/y-l-g/pogo/commit/87feb7e884dbef2ae5d6a42eae6cf04c482b49ae"
        },
        "date": 1764335249781,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkAllocate",
            "value": 84.58,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "14627428 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - ns/op",
            "value": 84.58,
            "unit": "ns/op",
            "extra": "14627428 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "14627428 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "14627428 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel",
            "value": 113.3,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "10350133 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - ns/op",
            "value": 113.3,
            "unit": "ns/op",
            "extra": "10350133 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "10350133 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "10350133 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt",
            "value": 54.5,
            "unit": "ns/op\t75152.15 MB/s\t       0 B/op\t       0 allocs/op",
            "extra": "21596816 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - ns/op",
            "value": 54.5,
            "unit": "ns/op",
            "extra": "21596816 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - MB/s",
            "value": 75152.15,
            "unit": "MB/s",
            "extra": "21596816 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "21596816 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "21596816 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON",
            "value": 3607,
            "unit": "ns/op\t    1360 B/op\t      36 allocs/op",
            "extra": "321800 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - ns/op",
            "value": 3607,
            "unit": "ns/op",
            "extra": "321800 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - B/op",
            "value": 1360,
            "unit": "B/op",
            "extra": "321800 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - allocs/op",
            "value": 36,
            "unit": "allocs/op",
            "extra": "321800 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack",
            "value": 694.7,
            "unit": "ns/op\t     192 B/op\t       1 allocs/op",
            "extra": "1725224 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - ns/op",
            "value": 694.7,
            "unit": "ns/op",
            "extra": "1725224 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - B/op",
            "value": 192,
            "unit": "B/op",
            "extra": "1725224 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1725224 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation",
            "value": 402.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "2953204 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - ns/op",
            "value": 402.4,
            "unit": "ns/op",
            "extra": "2953204 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "2953204 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "2953204 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus",
            "value": 657.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1589928 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - ns/op",
            "value": 657.6,
            "unit": "ns/op",
            "extra": "1589928 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1589928 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1589928 times\n4 procs"
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
          "id": "d6e244b5177fec32ae82d4907a230ebec0f771cf",
          "message": "wip",
          "timestamp": "2025-11-28T14:39:51+01:00",
          "tree_id": "0090163b7effb41c4c7092622e50b2f42d7ae4bc",
          "url": "https://github.com/y-l-g/pogo/commit/d6e244b5177fec32ae82d4907a230ebec0f771cf"
        },
        "date": 1764337241900,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkAllocate",
            "value": 91.37,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "12747120 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - ns/op",
            "value": 91.37,
            "unit": "ns/op",
            "extra": "12747120 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "12747120 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "12747120 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel",
            "value": 112.6,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "10656409 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - ns/op",
            "value": 112.6,
            "unit": "ns/op",
            "extra": "10656409 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "10656409 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "10656409 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt",
            "value": 54.79,
            "unit": "ns/op\t74757.29 MB/s\t       0 B/op\t       0 allocs/op",
            "extra": "21866899 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - ns/op",
            "value": 54.79,
            "unit": "ns/op",
            "extra": "21866899 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - MB/s",
            "value": 74757.29,
            "unit": "MB/s",
            "extra": "21866899 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "21866899 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "21866899 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON",
            "value": 3575,
            "unit": "ns/op\t    1360 B/op\t      36 allocs/op",
            "extra": "324723 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - ns/op",
            "value": 3575,
            "unit": "ns/op",
            "extra": "324723 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - B/op",
            "value": 1360,
            "unit": "B/op",
            "extra": "324723 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - allocs/op",
            "value": 36,
            "unit": "allocs/op",
            "extra": "324723 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack",
            "value": 701.6,
            "unit": "ns/op\t     192 B/op\t       1 allocs/op",
            "extra": "1708754 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - ns/op",
            "value": 701.6,
            "unit": "ns/op",
            "extra": "1708754 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - B/op",
            "value": 192,
            "unit": "B/op",
            "extra": "1708754 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1708754 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation",
            "value": 406.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "3001299 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - ns/op",
            "value": 406.6,
            "unit": "ns/op",
            "extra": "3001299 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "3001299 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "3001299 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus",
            "value": 637.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1875070 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - ns/op",
            "value": 637.2,
            "unit": "ns/op",
            "extra": "1875070 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1875070 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1875070 times\n4 procs"
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
          "id": "804023cfac8a980370a3be0b76f95cb9b47f4c55",
          "message": "wip",
          "timestamp": "2025-11-28T16:46:38+01:00",
          "tree_id": "d90f603e7b148f54f4e60e0cc503c5ec2143f038",
          "url": "https://github.com/y-l-g/pogo/commit/804023cfac8a980370a3be0b76f95cb9b47f4c55"
        },
        "date": 1764344851427,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkAllocate",
            "value": 86.48,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "11914402 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - ns/op",
            "value": 86.48,
            "unit": "ns/op",
            "extra": "11914402 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "11914402 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "11914402 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel",
            "value": 112.4,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "10548462 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - ns/op",
            "value": 112.4,
            "unit": "ns/op",
            "extra": "10548462 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "10548462 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "10548462 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt",
            "value": 98.43,
            "unit": "ns/op\t41613.66 MB/s\t       0 B/op\t       0 allocs/op",
            "extra": "21920986 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - ns/op",
            "value": 98.43,
            "unit": "ns/op",
            "extra": "21920986 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - MB/s",
            "value": 41613.66,
            "unit": "MB/s",
            "extra": "21920986 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "21920986 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "21920986 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON",
            "value": 3582,
            "unit": "ns/op\t    1360 B/op\t      36 allocs/op",
            "extra": "324351 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - ns/op",
            "value": 3582,
            "unit": "ns/op",
            "extra": "324351 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - B/op",
            "value": 1360,
            "unit": "B/op",
            "extra": "324351 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - allocs/op",
            "value": 36,
            "unit": "allocs/op",
            "extra": "324351 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack",
            "value": 694.3,
            "unit": "ns/op\t     192 B/op\t       1 allocs/op",
            "extra": "1718986 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - ns/op",
            "value": 694.3,
            "unit": "ns/op",
            "extra": "1718986 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - B/op",
            "value": 192,
            "unit": "B/op",
            "extra": "1718986 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1718986 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation",
            "value": 404.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "2935398 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - ns/op",
            "value": 404.6,
            "unit": "ns/op",
            "extra": "2935398 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "2935398 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "2935398 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus",
            "value": 637,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1884559 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - ns/op",
            "value": 637,
            "unit": "ns/op",
            "extra": "1884559 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1884559 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1884559 times\n4 procs"
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
          "id": "eb4dc95d5f0c73ac013aa112e5b79dbb835a022d",
          "message": "wip",
          "timestamp": "2025-11-28T17:19:15+01:00",
          "tree_id": "7423a68e92ef82caa059143dcedc3545cf8aa1ab",
          "url": "https://github.com/y-l-g/pogo/commit/eb4dc95d5f0c73ac013aa112e5b79dbb835a022d"
        },
        "date": 1764346803769,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkAllocate",
            "value": 88.44,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "13097221 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - ns/op",
            "value": 88.44,
            "unit": "ns/op",
            "extra": "13097221 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "13097221 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "13097221 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel",
            "value": 118.7,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "10278614 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - ns/op",
            "value": 118.7,
            "unit": "ns/op",
            "extra": "10278614 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "10278614 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "10278614 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt",
            "value": 54.45,
            "unit": "ns/op\t75231.39 MB/s\t       0 B/op\t       0 allocs/op",
            "extra": "21943483 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - ns/op",
            "value": 54.45,
            "unit": "ns/op",
            "extra": "21943483 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - MB/s",
            "value": 75231.39,
            "unit": "MB/s",
            "extra": "21943483 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "21943483 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "21943483 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON",
            "value": 3737,
            "unit": "ns/op\t    1360 B/op\t      36 allocs/op",
            "extra": "314284 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - ns/op",
            "value": 3737,
            "unit": "ns/op",
            "extra": "314284 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - B/op",
            "value": 1360,
            "unit": "B/op",
            "extra": "314284 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - allocs/op",
            "value": 36,
            "unit": "allocs/op",
            "extra": "314284 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack",
            "value": 700,
            "unit": "ns/op\t     192 B/op\t       1 allocs/op",
            "extra": "1705618 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - ns/op",
            "value": 700,
            "unit": "ns/op",
            "extra": "1705618 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - B/op",
            "value": 192,
            "unit": "B/op",
            "extra": "1705618 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1705618 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation",
            "value": 403.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "2999221 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - ns/op",
            "value": 403.8,
            "unit": "ns/op",
            "extra": "2999221 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "2999221 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "2999221 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus",
            "value": 661.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1805014 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - ns/op",
            "value": 661.6,
            "unit": "ns/op",
            "extra": "1805014 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1805014 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1805014 times\n4 procs"
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
          "id": "3717572066f898b4c748fed02ab9d51ed6033f61",
          "message": "wip",
          "timestamp": "2025-11-28T17:51:53+01:00",
          "tree_id": "02cf1624c0c65d65aea6d41127859f9b1eb808eb",
          "url": "https://github.com/y-l-g/pogo/commit/3717572066f898b4c748fed02ab9d51ed6033f61"
        },
        "date": 1764348768863,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkAllocate",
            "value": 77.44,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "14623669 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - ns/op",
            "value": 77.44,
            "unit": "ns/op",
            "extra": "14623669 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "14623669 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "14623669 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel",
            "value": 111.3,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "9072380 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - ns/op",
            "value": 111.3,
            "unit": "ns/op",
            "extra": "9072380 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "9072380 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "9072380 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt",
            "value": 54.54,
            "unit": "ns/op\t75106.75 MB/s\t       0 B/op\t       0 allocs/op",
            "extra": "21964368 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - ns/op",
            "value": 54.54,
            "unit": "ns/op",
            "extra": "21964368 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - MB/s",
            "value": 75106.75,
            "unit": "MB/s",
            "extra": "21964368 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "21964368 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "21964368 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON",
            "value": 3576,
            "unit": "ns/op\t    1360 B/op\t      36 allocs/op",
            "extra": "329565 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - ns/op",
            "value": 3576,
            "unit": "ns/op",
            "extra": "329565 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - B/op",
            "value": 1360,
            "unit": "B/op",
            "extra": "329565 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - allocs/op",
            "value": 36,
            "unit": "allocs/op",
            "extra": "329565 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack",
            "value": 704.3,
            "unit": "ns/op\t     192 B/op\t       1 allocs/op",
            "extra": "1725016 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - ns/op",
            "value": 704.3,
            "unit": "ns/op",
            "extra": "1725016 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - B/op",
            "value": 192,
            "unit": "B/op",
            "extra": "1725016 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1725016 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation",
            "value": 406.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "2954631 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - ns/op",
            "value": 406.1,
            "unit": "ns/op",
            "extra": "2954631 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "2954631 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "2954631 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus",
            "value": 639.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1866656 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - ns/op",
            "value": 639.2,
            "unit": "ns/op",
            "extra": "1866656 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1866656 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1866656 times\n4 procs"
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
          "id": "6662fb61d549f3a5aba8270fc9339da77376e8c2",
          "message": "wip",
          "timestamp": "2025-11-28T20:23:46+01:00",
          "tree_id": "19d624a50090154c2cd6ae7d5840c04801c740f2",
          "url": "https://github.com/y-l-g/pogo/commit/6662fb61d549f3a5aba8270fc9339da77376e8c2"
        },
        "date": 1764357880926,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkAllocate",
            "value": 83.21,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "13948219 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - ns/op",
            "value": 83.21,
            "unit": "ns/op",
            "extra": "13948219 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "13948219 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "13948219 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel",
            "value": 111.3,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "10413026 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - ns/op",
            "value": 111.3,
            "unit": "ns/op",
            "extra": "10413026 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "10413026 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "10413026 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt",
            "value": 54.62,
            "unit": "ns/op\t74995.40 MB/s\t       0 B/op\t       0 allocs/op",
            "extra": "21912878 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - ns/op",
            "value": 54.62,
            "unit": "ns/op",
            "extra": "21912878 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - MB/s",
            "value": 74995.4,
            "unit": "MB/s",
            "extra": "21912878 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "21912878 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "21912878 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON",
            "value": 3576,
            "unit": "ns/op\t    1360 B/op\t      36 allocs/op",
            "extra": "332318 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - ns/op",
            "value": 3576,
            "unit": "ns/op",
            "extra": "332318 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - B/op",
            "value": 1360,
            "unit": "B/op",
            "extra": "332318 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - allocs/op",
            "value": 36,
            "unit": "allocs/op",
            "extra": "332318 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack",
            "value": 702.3,
            "unit": "ns/op\t     192 B/op\t       1 allocs/op",
            "extra": "1681074 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - ns/op",
            "value": 702.3,
            "unit": "ns/op",
            "extra": "1681074 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - B/op",
            "value": 192,
            "unit": "B/op",
            "extra": "1681074 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1681074 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation",
            "value": 401.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "2949888 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - ns/op",
            "value": 401.9,
            "unit": "ns/op",
            "extra": "2949888 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "2949888 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "2949888 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus",
            "value": 641.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1867459 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - ns/op",
            "value": 641.8,
            "unit": "ns/op",
            "extra": "1867459 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1867459 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1867459 times\n4 procs"
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
          "id": "4b85e27610cccce4c7291930ee2c2eb2e30a9661",
          "message": "wip",
          "timestamp": "2025-12-01T22:28:43+01:00",
          "tree_id": "6fa41562fd19d2d5e98e2b3639388a6bbf8f2f0b",
          "url": "https://github.com/y-l-g/pogo/commit/4b85e27610cccce4c7291930ee2c2eb2e30a9661"
        },
        "date": 1764624583823,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkAllocate",
            "value": 83.71,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "14309665 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - ns/op",
            "value": 83.71,
            "unit": "ns/op",
            "extra": "14309665 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "14309665 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "14309665 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel",
            "value": 114.8,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "10374955 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - ns/op",
            "value": 114.8,
            "unit": "ns/op",
            "extra": "10374955 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "10374955 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "10374955 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt",
            "value": 54.97,
            "unit": "ns/op\t74511.42 MB/s\t       0 B/op\t       0 allocs/op",
            "extra": "21601800 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - ns/op",
            "value": 54.97,
            "unit": "ns/op",
            "extra": "21601800 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - MB/s",
            "value": 74511.42,
            "unit": "MB/s",
            "extra": "21601800 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "21601800 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "21601800 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON",
            "value": 3554,
            "unit": "ns/op\t    1360 B/op\t      36 allocs/op",
            "extra": "302292 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - ns/op",
            "value": 3554,
            "unit": "ns/op",
            "extra": "302292 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - B/op",
            "value": 1360,
            "unit": "B/op",
            "extra": "302292 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - allocs/op",
            "value": 36,
            "unit": "allocs/op",
            "extra": "302292 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack",
            "value": 701.3,
            "unit": "ns/op\t     192 B/op\t       1 allocs/op",
            "extra": "1723282 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - ns/op",
            "value": 701.3,
            "unit": "ns/op",
            "extra": "1723282 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - B/op",
            "value": 192,
            "unit": "B/op",
            "extra": "1723282 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1723282 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation",
            "value": 398.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "2987059 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - ns/op",
            "value": 398.1,
            "unit": "ns/op",
            "extra": "2987059 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "2987059 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "2987059 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus",
            "value": 653,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1830835 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - ns/op",
            "value": 653,
            "unit": "ns/op",
            "extra": "1830835 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1830835 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1830835 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "116063953+y-l-g@users.noreply.github.com",
            "name": "YL",
            "username": "y-l-g"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "19bcfd15ba03f23c8f7aca64bf42d86ba861aef5",
          "message": "Merge pull request #1 from y-l-g/renovate/configure\n\nConfigure Renovate",
          "timestamp": "2026-01-03T23:27:39+01:00",
          "tree_id": "b4b33601f37ac08e68379b8b69b29118afa0b865",
          "url": "https://github.com/y-l-g/pogo/commit/19bcfd15ba03f23c8f7aca64bf42d86ba861aef5"
        },
        "date": 1767479313980,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkAllocate",
            "value": 88.01,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "14338747 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - ns/op",
            "value": 88.01,
            "unit": "ns/op",
            "extra": "14338747 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "14338747 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "14338747 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel",
            "value": 116.4,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "10619046 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - ns/op",
            "value": 116.4,
            "unit": "ns/op",
            "extra": "10619046 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "10619046 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "10619046 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt",
            "value": 54.54,
            "unit": "ns/op\t75104.82 MB/s\t       0 B/op\t       0 allocs/op",
            "extra": "21666931 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - ns/op",
            "value": 54.54,
            "unit": "ns/op",
            "extra": "21666931 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - MB/s",
            "value": 75104.82,
            "unit": "MB/s",
            "extra": "21666931 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "21666931 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "21666931 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON",
            "value": 3556,
            "unit": "ns/op\t    1360 B/op\t      36 allocs/op",
            "extra": "331129 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - ns/op",
            "value": 3556,
            "unit": "ns/op",
            "extra": "331129 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - B/op",
            "value": 1360,
            "unit": "B/op",
            "extra": "331129 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - allocs/op",
            "value": 36,
            "unit": "allocs/op",
            "extra": "331129 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack",
            "value": 712.4,
            "unit": "ns/op\t     192 B/op\t       1 allocs/op",
            "extra": "1728432 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - ns/op",
            "value": 712.4,
            "unit": "ns/op",
            "extra": "1728432 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - B/op",
            "value": 192,
            "unit": "B/op",
            "extra": "1728432 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1728432 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation",
            "value": 404.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "3030007 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - ns/op",
            "value": 404.2,
            "unit": "ns/op",
            "extra": "3030007 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "3030007 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "3030007 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus",
            "value": 643.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1870026 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - ns/op",
            "value": 643.8,
            "unit": "ns/op",
            "extra": "1870026 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1870026 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1870026 times\n4 procs"
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
          "id": "a3a89d212143c3b60d3b82c277123ceb7769eea6",
          "message": "wip",
          "timestamp": "2026-01-03T23:44:40+01:00",
          "tree_id": "e4f4eb60b68682c5c5be002d1f1a4b38fe36c07f",
          "url": "https://github.com/y-l-g/pogo/commit/a3a89d212143c3b60d3b82c277123ceb7769eea6"
        },
        "date": 1767480314025,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkAllocate",
            "value": 77.01,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "15473424 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - ns/op",
            "value": 77.01,
            "unit": "ns/op",
            "extra": "15473424 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "15473424 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "15473424 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel",
            "value": 116.8,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "10493743 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - ns/op",
            "value": 116.8,
            "unit": "ns/op",
            "extra": "10493743 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "10493743 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "10493743 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt",
            "value": 98.36,
            "unit": "ns/op\t41643.51 MB/s\t       0 B/op\t       0 allocs/op",
            "extra": "21957652 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - ns/op",
            "value": 98.36,
            "unit": "ns/op",
            "extra": "21957652 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - MB/s",
            "value": 41643.51,
            "unit": "MB/s",
            "extra": "21957652 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "21957652 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "21957652 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON",
            "value": 3588,
            "unit": "ns/op\t    1360 B/op\t      36 allocs/op",
            "extra": "326046 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - ns/op",
            "value": 3588,
            "unit": "ns/op",
            "extra": "326046 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - B/op",
            "value": 1360,
            "unit": "B/op",
            "extra": "326046 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - allocs/op",
            "value": 36,
            "unit": "allocs/op",
            "extra": "326046 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack",
            "value": 698.7,
            "unit": "ns/op\t     192 B/op\t       1 allocs/op",
            "extra": "1713024 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - ns/op",
            "value": 698.7,
            "unit": "ns/op",
            "extra": "1713024 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - B/op",
            "value": 192,
            "unit": "B/op",
            "extra": "1713024 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1713024 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation",
            "value": 405.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "3006007 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - ns/op",
            "value": 405.7,
            "unit": "ns/op",
            "extra": "3006007 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "3006007 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "3006007 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus",
            "value": 693.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1866498 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - ns/op",
            "value": 693.3,
            "unit": "ns/op",
            "extra": "1866498 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1866498 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1866498 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "116063953+y-l-g@users.noreply.github.com",
            "name": "YL",
            "username": "y-l-g"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "622565c2f828df5157577f60ebdd5cad8bc1987e",
          "message": "Merge pull request #2 from y-l-g/renovate/github.com-dunglas-frankenphp-1.x\n\nUpdate module github.com/dunglas/frankenphp to v1.11.1",
          "timestamp": "2026-01-03T23:51:44+01:00",
          "tree_id": "4ced15c40485083e9aca26646a0cf1c5b29d69ae",
          "url": "https://github.com/y-l-g/pogo/commit/622565c2f828df5157577f60ebdd5cad8bc1987e"
        },
        "date": 1767480760049,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkAllocate",
            "value": 76.99,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "13400116 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - ns/op",
            "value": 76.99,
            "unit": "ns/op",
            "extra": "13400116 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "13400116 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "13400116 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel",
            "value": 113.2,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "10421919 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - ns/op",
            "value": 113.2,
            "unit": "ns/op",
            "extra": "10421919 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "10421919 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "10421919 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt",
            "value": 54.94,
            "unit": "ns/op\t74552.99 MB/s\t       0 B/op\t       0 allocs/op",
            "extra": "21920994 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - ns/op",
            "value": 54.94,
            "unit": "ns/op",
            "extra": "21920994 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - MB/s",
            "value": 74552.99,
            "unit": "MB/s",
            "extra": "21920994 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "21920994 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "21920994 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON",
            "value": 3562,
            "unit": "ns/op\t    1360 B/op\t      36 allocs/op",
            "extra": "329416 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - ns/op",
            "value": 3562,
            "unit": "ns/op",
            "extra": "329416 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - B/op",
            "value": 1360,
            "unit": "B/op",
            "extra": "329416 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - allocs/op",
            "value": 36,
            "unit": "allocs/op",
            "extra": "329416 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack",
            "value": 695.8,
            "unit": "ns/op\t     192 B/op\t       1 allocs/op",
            "extra": "1720270 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - ns/op",
            "value": 695.8,
            "unit": "ns/op",
            "extra": "1720270 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - B/op",
            "value": 192,
            "unit": "B/op",
            "extra": "1720270 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1720270 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation",
            "value": 404.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "3020060 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - ns/op",
            "value": 404.5,
            "unit": "ns/op",
            "extra": "3020060 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "3020060 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "3020060 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus",
            "value": 646.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1827946 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - ns/op",
            "value": 646.2,
            "unit": "ns/op",
            "extra": "1827946 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1827946 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1827946 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "116063953+y-l-g@users.noreply.github.com",
            "name": "YL",
            "username": "y-l-g"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "fdd5486eb324baf3180510aa830fad7cb7858706",
          "message": "Merge pull request #5 from y-l-g/renovate/actions-checkout-6.x\n\nUpdate actions/checkout action to v6",
          "timestamp": "2026-01-03T23:52:36+01:00",
          "tree_id": "78682f4a18c2d15d7519263b8c65c1feb07a8ed5",
          "url": "https://github.com/y-l-g/pogo/commit/fdd5486eb324baf3180510aa830fad7cb7858706"
        },
        "date": 1767480782439,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkAllocate",
            "value": 76.75,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "15474063 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - ns/op",
            "value": 76.75,
            "unit": "ns/op",
            "extra": "15474063 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "15474063 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "15474063 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel",
            "value": 112.2,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "10769353 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - ns/op",
            "value": 112.2,
            "unit": "ns/op",
            "extra": "10769353 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "10769353 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "10769353 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt",
            "value": 54.47,
            "unit": "ns/op\t75194.60 MB/s\t       0 B/op\t       0 allocs/op",
            "extra": "21642090 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - ns/op",
            "value": 54.47,
            "unit": "ns/op",
            "extra": "21642090 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - MB/s",
            "value": 75194.6,
            "unit": "MB/s",
            "extra": "21642090 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "21642090 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "21642090 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON",
            "value": 3547,
            "unit": "ns/op\t    1360 B/op\t      36 allocs/op",
            "extra": "327894 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - ns/op",
            "value": 3547,
            "unit": "ns/op",
            "extra": "327894 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - B/op",
            "value": 1360,
            "unit": "B/op",
            "extra": "327894 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - allocs/op",
            "value": 36,
            "unit": "allocs/op",
            "extra": "327894 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack",
            "value": 690.8,
            "unit": "ns/op\t     192 B/op\t       1 allocs/op",
            "extra": "1728595 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - ns/op",
            "value": 690.8,
            "unit": "ns/op",
            "extra": "1728595 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - B/op",
            "value": 192,
            "unit": "B/op",
            "extra": "1728595 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1728595 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation",
            "value": 399.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "2982830 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - ns/op",
            "value": 399.5,
            "unit": "ns/op",
            "extra": "2982830 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "2982830 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "2982830 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus",
            "value": 638.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1872806 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - ns/op",
            "value": 638.6,
            "unit": "ns/op",
            "extra": "1872806 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1872806 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1872806 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "116063953+y-l-g@users.noreply.github.com",
            "name": "YL",
            "username": "y-l-g"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "a63fa5a42ee71b72083d678956022e1f22d44346",
          "message": "Merge pull request #6 from y-l-g/renovate/actions-github-script-8.x\n\nUpdate actions/github-script action to v8",
          "timestamp": "2026-01-03T23:52:49+01:00",
          "tree_id": "d582229588082fba2af73f37a885216c6019b6d4",
          "url": "https://github.com/y-l-g/pogo/commit/a63fa5a42ee71b72083d678956022e1f22d44346"
        },
        "date": 1767480823355,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkAllocate",
            "value": 77.38,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "15319984 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - ns/op",
            "value": 77.38,
            "unit": "ns/op",
            "extra": "15319984 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "15319984 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "15319984 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel",
            "value": 114.6,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "10334850 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - ns/op",
            "value": 114.6,
            "unit": "ns/op",
            "extra": "10334850 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "10334850 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "10334850 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt",
            "value": 54.65,
            "unit": "ns/op\t74945.39 MB/s\t       0 B/op\t       0 allocs/op",
            "extra": "21946195 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - ns/op",
            "value": 54.65,
            "unit": "ns/op",
            "extra": "21946195 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - MB/s",
            "value": 74945.39,
            "unit": "MB/s",
            "extra": "21946195 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "21946195 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "21946195 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON",
            "value": 3589,
            "unit": "ns/op\t    1360 B/op\t      36 allocs/op",
            "extra": "332276 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - ns/op",
            "value": 3589,
            "unit": "ns/op",
            "extra": "332276 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - B/op",
            "value": 1360,
            "unit": "B/op",
            "extra": "332276 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - allocs/op",
            "value": 36,
            "unit": "allocs/op",
            "extra": "332276 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack",
            "value": 690.6,
            "unit": "ns/op\t     192 B/op\t       1 allocs/op",
            "extra": "1733133 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - ns/op",
            "value": 690.6,
            "unit": "ns/op",
            "extra": "1733133 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - B/op",
            "value": 192,
            "unit": "B/op",
            "extra": "1733133 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1733133 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation",
            "value": 403.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "2918409 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - ns/op",
            "value": 403.6,
            "unit": "ns/op",
            "extra": "2918409 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "2918409 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "2918409 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus",
            "value": 650.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1828120 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - ns/op",
            "value": 650.9,
            "unit": "ns/op",
            "extra": "1828120 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1828120 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1828120 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "116063953+y-l-g@users.noreply.github.com",
            "name": "YL",
            "username": "y-l-g"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "475313bf4ef31cee7295abfb3b9ed801a3feb275",
          "message": "Merge pull request #7 from y-l-g/renovate/major-github-artifact-actions\n\nUpdate GitHub Artifact Actions (major)",
          "timestamp": "2026-01-03T23:53:02+01:00",
          "tree_id": "c8e7c3fe2b6669ab46f93386c16b80e34201514e",
          "url": "https://github.com/y-l-g/pogo/commit/475313bf4ef31cee7295abfb3b9ed801a3feb275"
        },
        "date": 1767480907734,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkAllocate",
            "value": 76.41,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "15531290 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - ns/op",
            "value": 76.41,
            "unit": "ns/op",
            "extra": "15531290 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "15531290 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "15531290 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel",
            "value": 111.7,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "10667995 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - ns/op",
            "value": 111.7,
            "unit": "ns/op",
            "extra": "10667995 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "10667995 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "10667995 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt",
            "value": 54.63,
            "unit": "ns/op\t74971.02 MB/s\t       0 B/op\t       0 allocs/op",
            "extra": "21906166 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - ns/op",
            "value": 54.63,
            "unit": "ns/op",
            "extra": "21906166 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - MB/s",
            "value": 74971.02,
            "unit": "MB/s",
            "extra": "21906166 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "21906166 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "21906166 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON",
            "value": 3586,
            "unit": "ns/op\t    1360 B/op\t      36 allocs/op",
            "extra": "327411 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - ns/op",
            "value": 3586,
            "unit": "ns/op",
            "extra": "327411 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - B/op",
            "value": 1360,
            "unit": "B/op",
            "extra": "327411 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - allocs/op",
            "value": 36,
            "unit": "allocs/op",
            "extra": "327411 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack",
            "value": 692.2,
            "unit": "ns/op\t     192 B/op\t       1 allocs/op",
            "extra": "1723398 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - ns/op",
            "value": 692.2,
            "unit": "ns/op",
            "extra": "1723398 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - B/op",
            "value": 192,
            "unit": "B/op",
            "extra": "1723398 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1723398 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation",
            "value": 401,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "3007178 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - ns/op",
            "value": 401,
            "unit": "ns/op",
            "extra": "3007178 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "3007178 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "3007178 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus",
            "value": 638.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1874652 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - ns/op",
            "value": 638.5,
            "unit": "ns/op",
            "extra": "1874652 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1874652 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1874652 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "116063953+y-l-g@users.noreply.github.com",
            "name": "YL",
            "username": "y-l-g"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "dc24bb4b4876b532c9cade95b353884682fef58b",
          "message": "Merge pull request #8 from y-l-g/renovate/golangci-golangci-lint-action-9.x\n\nUpdate golangci/golangci-lint-action action to v9",
          "timestamp": "2026-01-03T23:53:13+01:00",
          "tree_id": "ee88c171180089d36f256abc124c4e959e19c9af",
          "url": "https://github.com/y-l-g/pogo/commit/dc24bb4b4876b532c9cade95b353884682fef58b"
        },
        "date": 1767480921354,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkAllocate",
            "value": 80.27,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "14777500 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - ns/op",
            "value": 80.27,
            "unit": "ns/op",
            "extra": "14777500 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "14777500 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocate - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "14777500 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel",
            "value": 129.1,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "8677668 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - ns/op",
            "value": 129.1,
            "unit": "ns/op",
            "extra": "8677668 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "8677668 times\n4 procs"
          },
          {
            "name": "BenchmarkAllocateParallel - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "8677668 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt",
            "value": 70.16,
            "unit": "ns/op\t58377.87 MB/s\t       0 B/op\t       0 allocs/op",
            "extra": "17073864 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - ns/op",
            "value": 70.16,
            "unit": "ns/op",
            "extra": "17073864 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - MB/s",
            "value": 58377.87,
            "unit": "MB/s",
            "extra": "17073864 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "17073864 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteAt - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "17073864 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON",
            "value": 3742,
            "unit": "ns/op\t    1360 B/op\t      36 allocs/op",
            "extra": "317715 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - ns/op",
            "value": 3742,
            "unit": "ns/op",
            "extra": "317715 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - B/op",
            "value": 1360,
            "unit": "B/op",
            "extra": "317715 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/JSON - allocs/op",
            "value": 36,
            "unit": "allocs/op",
            "extra": "317715 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack",
            "value": 674.5,
            "unit": "ns/op\t     192 B/op\t       1 allocs/op",
            "extra": "1779021 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - ns/op",
            "value": 674.5,
            "unit": "ns/op",
            "extra": "1779021 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - B/op",
            "value": 192,
            "unit": "B/op",
            "extra": "1779021 times\n4 procs"
          },
          {
            "name": "BenchmarkSerialization/MsgPack - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1779021 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation",
            "value": 424.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "2814187 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - ns/op",
            "value": 424.9,
            "unit": "ns/op",
            "extra": "2814187 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "2814187 times\n4 procs"
          },
          {
            "name": "BenchmarkHandleValidation - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "2814187 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus",
            "value": 706.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1740603 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - ns/op",
            "value": 706.3,
            "unit": "ns/op",
            "extra": "1740603 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1740603 times\n4 procs"
          },
          {
            "name": "BenchmarkInternalBus - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1740603 times\n4 procs"
          }
        ]
      }
    ]
  }
}