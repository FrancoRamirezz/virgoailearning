import Hero from "@/components/landing/Hero";
import Features from "@/components/landing/Features";
import HowItWorks from "@/components/landing/HowItWorks";
import Readiness from "@/components/landing/Readiness";
import MissionTeaser from "@/components/landing/MissionTeaser";
import BottomCTA from "@/components/landing/BottomCTA";

export default function Page() {
  return (
    <main className="flex min-h-screen flex-col gradient-bg">
        <Hero />
        <Features />
        <HowItWorks />
        <Readiness />
        <MissionTeaser />
        <BottomCTA />
    </main>
  );
}

