// components/landing/Features.tsx
import { 
  BookOpen, 
  Users, 
  Award, 
  Star, 
} from 'lucide-react';
export default function Features() {
  return (
    <section className="section ">
      <div className="section-inner">
        <div className="text-center max-w-2xl mx-auto">
          <h2>Everything You Need to Succeed in Your New Home</h2>
          <p className="mt-3 text-sm md:text-base text-slate-600">
            From citizenship prep to everyday English and life skills,
            VirgoLearning guides you step by step through the journey of
            living, working, and thriving in the United States.
          </p>
        </div>

        <div className="mt-10 grid gap-6 md:grid-cols-3">
          <div className="card">
            <div className="mb-4 h-9 w-9 rounded-2xl bg-[var(--color-navy)]/10" />
            <h3>Citizenship Test Prep</h3>
            <p className="mt-3 text-sm text-slate-600">
              Learn all the civics questions, practice interview-style answers,
              and track your readiness with realistic practice tests and simple
              explanations.
            </p>
          </div>

          <div className="card">
            <div className="mb-4 h-9 w-9 rounded-2xl bg-[var(--color-gold)]/15" />
            <h3>Immigrant-Focused English</h3>
            <p className="mt-3 text-sm text-slate-600">
              Build the English you actually need — for work, school,
              appointments, and daily conversations — with bilingual support for
              Spanish and Chinese speakers.
            </p>
          </div>

          <div className="card">
            <div className="mb-4 h-9 w-9 rounded-2xl bg-[var(--color-red)]/12" />
            <h3>Life Skills & Onboarding</h3>
            <p className="mt-3 text-sm text-slate-600">
              Practice real-life situations like calling a doctor, talking to
              your child’s teacher, or speaking with your landlord so you feel
              confident in everyday life.
            </p>
          </div>
        </div>
      </div>
    </section>
  );
}
