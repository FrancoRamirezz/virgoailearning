// components/landing/Readiness.tsx
export default function Readiness() {
  return (
    <section className="section" id="readiness">
      <div className="section-inner grid gap-8 md:grid-cols-2 md:items-center">
        <div>
          <h2>See Your Readiness at Every Step</h2>
          <p className="mt-3 text-sm md:text-base text-slate-600">
            Instead of guessing if you’re ready, VirgoLearning shows your
            progress clearly — for the test and for real life.
          </p>

          <p className="mt-4 text-sm text-slate-600">
            Your personal dashboard tracks how prepared you are in three key
            areas:
          </p>

          <ul className="mt-4 space-y-2 text-sm text-slate-700">
            <li>
              <span className="font-semibold">Citizenship Readiness</span> – based on civics modules and
              practice tests.
            </li>
            <li>
              <span className="font-semibold">English Readiness</span> – based on vocabulary, grammar, and
              speaking practice.
            </li>
            <li>
              <span className="font-semibold">Life Skills Readiness</span> – based on real-world scenarios
              like healthcare, school, housing, and money.
            </li>
          </ul>

          <p className="mt-4 text-sm text-slate-600">
            As you learn, your scores update automatically, and the platform
            suggests what to study next so you always know your next best step.
          </p>
        </div>

        {/* Visual mockup */}
        <div className="card-soft">
          <p className="text-xs font-medium uppercase tracking-wide text-slate-500 mb-4">
            Example readiness dashboard
          </p>

          <div className="space-y-4">
            {/* Progress bars */}
            {[
              { label: "Citizenship Readiness", value: 82 },
              { label: "English Readiness", value: 64 },
              { label: "Life Skills Readiness", value: 58 },
            ].map((item) => (
              <div key={item.label}>
                <div className="flex justify-between text-xs text-slate-600 mb-1">
                  <span>{item.label}</span>
                  <span className="font-semibold text-slate-800">
                    {item.value}%
                  </span>
                </div>
                <div className="h-2 w-full rounded-full bg-[var(--color-light-gray)]">
                  <div
                    className="h-2 rounded-full bg-[var(--color-navy)]"
                    style={{ width: `${item.value}%` }}
                  />
                </div>
              </div>
            ))}

            <div className="mt-6 rounded-2xl bg-white p-4 shadow-sm border border-[var(--color-light-gray)]">
              <p className="text-xs font-semibold text-slate-800">
                Suggested next step
              </p>
              <p className="mt-2 text-xs text-slate-600">
                Review: “Talking to your child&apos;s teacher” and retake the
                Life Skills quiz to boost your readiness.
              </p>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
