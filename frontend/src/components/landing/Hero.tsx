import { CheckCircle } from "lucide-react";

export default function Hero() {
  return (
    <section className="section">
      <div className="section-inner flex flex-col gap-12 md:flex-row md:items-center">
        
        {/* Left copy */}
        <div className="max-w-xl">
          
          {/* Pre-headline */}
          <span className="badge-pill mb-4">
            Based on official USCIS civics questions
          </span>

          <h1>
            Know When You’re Ready for U.S. Citizenship.
          </h1>

          <p className="mt-4 text-sm md:text-base text-slate-700">
            Practice real citizenship questions, build confidence in English,
            and track your readiness — so you walk into your interview knowing
            exactly where you stand.
          </p>

          <div className="mt-6 flex flex-wrap gap-3">
            <button className="btn-primary">
              Start Free Practice
            </button>
            <button className="btn-secondary">
              See How It Works
            </button>
          </div>

          {/* Trust signals (facts, not claims) */}
          <div className="mt-8 flex flex-wrap gap-4 text-xs md:text-sm text-slate-600">
            <div className="pill-metric">
              <span className="font-semibold text-slate-900">100</span>
              <span>official civics questions</span>
            </div>
            <div className="pill-metric">
              <span className="font-semibold text-slate-900">Bilingual</span>
              <span>simple English support</span>
            </div>
            <div className="pill-metric">
              <span className="font-semibold text-slate-900">Mobile</span>
              <span>learn anywhere</span>
            </div>
          </div>

          <p className="mt-3 text-[11px] text-slate-400">
            Free to start. No credit card required.
          </p>
        </div>

        {/* Right visual / outcome card */}
        <div className="flex-1">
          <div className="card-soft h-full min-h-[2090px] md:min-h-[240px] flex flex-col justify-between">
            
            <div>
              <p className="text-base md:text-lg font-bold text-slate-800">
                See your readiness grow as you practice.
              </p>
              <p className="mt-3 text-s text-slate-600">
                Your dashboard tracks civics, English, and real-life skills —
                showing what you’ve mastered and what to review next.
              </p>
            </div>

            <div className="mt-6 inline-flex items-center gap-3 self-start rounded-full bg-white px-4 py-2 shadow-sm border border-[var(--color-light-gray)]">
              <div className="w-8 h-8 bg-gradient-to-r from-red-100 to-blue-100 rounded-full flex items-center justify-center">
                <CheckCircle className="w-6 h-6 text-green-600" />
              </div>
              <div>
                <p className="text-xs font-semibold text-slate-900">
                  Clear next steps
                </p>
                <p className="text-[11px] text-slate-500">
                  based on your progress
                </p>
              </div>
            </div>

          </div>
        </div>

      </div>
    </section>
  );
}


// // components/landing/Hero.tsx
// import { CheckCircle} from 'lucide-react';

// export default function Hero() {
//   return (
//     <section className="section">
//       <div className="section-inner flex flex-col gap-10 md:flex-row md:items-center">
//         {/* Left copy */}
//         <div className="max-w-xl">
//           <span className="badge-pill mb-4">
//             Immigrant Learning & Onboarding Platform
//           </span>

//           <h1>Know When You’re Ready for U.S. Citizenship.</h1>

//           <p className="mt-4 text-sm md:text-base text-slate-700">
//             Practice real civics questions, improve your English, and track your readiness — so you stop guessing and walk into your interview with confidence.
//           </p>

//           <div className="mt-6 flex flex-wrap gap-3">
//             <button className="btn-primary">Get Started Free</button>
//             <button className="btn-secondary">Explore How It Works</button>
//           </div>

//           <div className="mt-8 flex flex-wrap gap-4 text-xs md:text-sm text-slate-600">
//             <div className="pill-metric">
//               <span className="font-semibold text-slate-900">98%</span>
//               <span>practice-test pass rate*</span>
//             </div>
//             <div className="pill-metric">
//               <span className="font-semibold text-slate-900">10,000+</span>
//               <span>lessons completed</span>
//             </div>
//             <div className="pill-metric">
//               <span className="font-semibold text-slate-900">4.9★</span>
//               <span>learner satisfaction</span>
//             </div>
//           </div>

//           <p className="mt-3 text-[11px] text-slate-400">
//             *Practice-test pass rate based on learners who completed all
//             recommended modules.
//           </p>
//         </div>

//         {/* Right visual placeholder */}
//         <div className="flex-1">
//           <div className="card-soft h-full min-h-[220px] md:min-h-[260px] flex flex-col justify-between">
//             <div>
//               <p className="text-sm font-medium text-slate-800">
//                 “I finally felt ready to book my citizenship interview.”
//               </p>
//               <p className="mt-3 text-xs text-slate-600">
//                 Practice civics questions, everyday English, and real-life
//                 scenarios in one dashboard. Visualize your readiness before you
//                 take the next step.
//               </p>
//             </div>
//             <div className="mt-6 inline-flex items-center gap-3 self-start rounded-full bg-white px-4 py-2 shadow-sm border border-[var(--color-light-gray)]">
//               <div className="w-8 h-8 bg-gradient-to-r from-red-100 to-blue-100 rounded-full flex items-center justify-center">
//                   <CheckCircle className="w-6 h-6 text-green-600" />
//                 </div>
//               <div>
//                 <p className="text-xs font-semibold text-slate-900">
//                   Certificate-ready in ~6 weeks
//                 </p>
//                 <p className="text-[11px] text-slate-500">
//                   with a consistent study plan
//                 </p>
//               </div>
//             </div>
//           </div>
//         </div>
//       </div>
//     </section>
//   );
// }
