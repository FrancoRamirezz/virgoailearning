// components/landing/MissionTeaser.tsx
import Link from "next/link";

export default function MissionTeaser() {
  return (
    <section className="section">
      <div className="section-inner grid gap-8 md:grid-cols-[1.3fr,1fr] md:items-center">
        <div>
          <h2>Built for Immigrants, With Immigrants</h2>
          <p className="mt-3 text-sm md:text-base text-slate-600">
            VirgoLearning was created with one goal: to make the journey into
            American life less confusing and more empowering.
          </p>
          <p className="mt-3 text-sm text-slate-600">
            We focus on Spanish and Chinese speakers first, combining
            expert-designed lessons, clear language support, and AI-powered
            explanations so you can learn at your own pace — on your phone, at
            home, or between work shifts.
          </p>

          <Link
            href="/about"
            className="mt-5 inline-flex items-center text-sm font-medium text-[var(--color-navy)] underline-offset-2 hover:underline"
          >
            Learn more about our mission →
          </Link>
        </div>

        <div className="card">
          <p className="text-sm font-semibold text-slate-800">
            “I used VirgoLearning after long shifts. The short lessons and
            simple explanations made it possible to keep going every day.”
          </p>
          <p className="mt-3 text-xs text-slate-500">– Example learner story</p>
        </div>
      </div>
    </section>
  );
}
