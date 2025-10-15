import React, { useState } from "react";
import { SafeAreaView, View, Text } from "react-native";
import { useSafeAreaInsets } from "react-native-safe-area-context";
import { useRouter } from "expo-router";
import VerificationJobList from "@/components/verifiers/VerificationJobList";
import DocumentReviewModal from "@/components/verifiers/DocumentReviewModal";
import { useKYCStore } from "@/store/kycStore";

const mockJobs = [
  { 
    id: "1", 
    requester: "Brandon", 
    doc: "Identity Card", 
    status: "pending",
    documentImage: "https://example.com/id1.jpg"
  },
  { 
    id: "2", 
    requester: "Brandon", 
    doc: "Identity Card", 
    status: "pending",
    documentImage: "https://example.com/id2.jpg"
  },
  { 
    id: "3", 
    requester: "Brandon", 
    doc: "Identity Card", 
    status: "pending",
    documentImage: "https://example.com/id3.jpg"
  },
  { 
    id: "4", 
    requester: "Brandon", 
    doc: "Identity Card", 
    status: "pending",
    documentImage: "https://example.com/id4.jpg"
  },
];

export default function VerifierDashboard() {
  const insets = useSafeAreaInsets();
  const router = useRouter();
  const setVerificationStatus = useKYCStore((state) => state.setVerificationStatus);
  
  const [jobs, setJobs] = useState(mockJobs);
  const [selectedJob, setSelectedJob] = useState(null);

  const handleApprove = (job) => {
    setVerificationStatus("verified");
    setJobs((prev) => prev.filter((j) => j.id !== job.id));
    setSelectedJob(null);
    
    if (jobs.length <= 1) {
      router.replace("/(tabs)/home");
    }
  };

  const handleReject = (job) => {
    setVerificationStatus("rejected");
    setJobs((prev) => prev.filter((j) => j.id !== job.id));
    setSelectedJob(null);
    
    if (jobs.length <= 1) {
      router.replace("/(tabs)/home");
    }
  };

  return (
    <SafeAreaView
      className="flex-1 bg-gray-50 px-5"
      style={{ paddingTop: insets.top }}
    >
      <VerificationJobList
        jobs={jobs}
        onJobPress={setSelectedJob}
      />
      
      <DocumentReviewModal
        visible={!!selectedJob}
        job={selectedJob}
        onClose={() => setSelectedJob(null)}
        onApprove={handleApprove}
        onReject={handleReject}
      />
    </SafeAreaView>
  );
}