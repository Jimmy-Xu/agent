package main

import (
	"errors"
	pb "github.com/kata-containers/agent/protocols/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
)

func getPodBuddyStats() (*pb.PodBuddyStats, error) {
}

func getPodCpuStats() ([]*pb.PodCpuStats, error) {
}
