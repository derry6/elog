#!/usr/bin/env bash

mkdir -p logs
name=$1
profile=$2
run_bench() {
	opts="-v -benchmem"
	if [ -n "$name" ]; then
	  opts="$opts -bench $name"
	else
	  opts="$opts -bench ."
	fi
	opts="$opts -benchtime=10s"
	if [ -n "$profile" ] ; then
	  opts="$opts -cpuprofile=${profile}.out"
	fi

	echo "Running: go test $opts"
	go test $opts

	if [ -n "$profile" ] ; then
	  echo "Generating pprof svg ..."
	  go tool pprof -svg -output ${profile}.svg ${profile}.out
	fi
}

case "$1" in
  clean)
    rm -fr logs *.test *.svg *.out
    ;;
  *)
	  run_bench
    ;;
esac

#rm -f *.elog *.svg *.out *.test
#go test -v -bench=. -benchtime=15s -cpuprofile=bench-pprof.out
#go tool pprof -svg -output bench-pprof.svg bench-pprof.out
