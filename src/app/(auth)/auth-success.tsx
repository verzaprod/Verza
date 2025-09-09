import Success from "@/components/ui/Success";

export default function AuthSuccess() {
  return (
    <Success 
      title="Success!"
      tagline="Your wallet has been created successfully. Let’s verify your identity to unlock all features."
      buttonText="Start KYC Verification"
      redirectType="auth"
    />
  )
}
