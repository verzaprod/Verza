import Success from "@/components/ui/Success";

export default function KYCSuccess() {
  return (
    <Success 
      redirectType="kyc"
      title="Credential Issued!"
      tagline="Your identity has been successfully verified and your digital credential is ready to use."
      buttonText="Go to Wallet"
    />
  )
}
