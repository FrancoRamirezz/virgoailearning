// components/landing/HowItWorks.tsx
const steps = [
  {
    title: "1. Choose Your Path",
    body: "Start with citizenship prep, English for work, or everyday life skills. The platform recommends lessons based on a short placement quiz.",
  },
  {
    title: "2. Learn with Short, Bilingual Lessons",
    body: "Study bite-sized modules with simple English, native-language support, practice questions, and real-world examples.",
  },
  {
    title: "3. See Your Readiness Grow",
    body: "Watch your Citizenship, English, and Life Skills readiness scores improve as you complete lessons and pass practice tests.",
  },
];

export default function HowItWorks() {
  return (
    <section className="section ">
      <div className="section-inner">
        <div className="text-center max-w-2xl mx-auto">
          <h2>How VirgoLearning Works</h2>
          <p className="mt-3 text-sm md:text-base text-slate-600">
            Simple, structured learning designed for busy immigrant families and
            workers.
          </p>
        </div>

        <div className="mt-10 grid gap-6 md:grid-cols-3">
          {steps.map((step) => (
            <div key={step.title} className="card-soft">
              <h3 className="text-base md:text-lg">{step.title}</h3>
              <p className="mt-3 text-sm text-slate-600">{step.body}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
