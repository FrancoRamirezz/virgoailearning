// components/landing/BottomCTA.tsx
import { Award } from 'lucide-react';
export default function BottomCTA() {
  return (
    <section className="section ">
      <div className='section-inner'>
        <div className="max-w-3xl mx-auto bg-gradient-to-r from-red-600 to-blue-600 rounded-2xl p-12 text-white text-center">
          <div className="w-16 h-16 bg-white/20 rounded-full flex items-center justify-center mx-auto mb-6">
            <Award className="w-8 h-8 text-white" />
          </div>
          <h2 className="text-4xl font-bold mb-6">Start Your Journey Today</h2>
          <p className="text-xl mb-8 text-purple-100 leading-relaxed">
            Join other immigrants who are studying smarter — not alone. Get clear lessons, real-life practice, and a simple way to see how ready you are for citizenship, work, and daily life in the U.S.
          </p>
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <button className="bg-white text-blue-600 font-bold py-3 px-8 rounded-xl transition-all duration-300 hover:bg-gray-100 transform hover:scale-105 shadow-lg hover:shadow-xl">
              Begin Learning Now
            </button>
            
          </div>
        </div>
      </div>
      
    </section>
    
  );
}
    // <section className="section">
    //   <div className="section-inner">
    //     <div className="rounded-3xl bg-gradient-to-r from-[var(--color-navy)] via-[var(--color-muted-blue)] to-[var(--color-red)] text-white px-6 py-10 md:px-10 md:py-12 shadow-lg">
    //       <div className="max-w-2xl">
    //         <h2 className="text-white">
    //           Ready to Start Your Journey?
    //         </h2>
    //         <p className="mt-3 text-sm md:text-base text-slate-100">
    //           Join other immigrants who are studying smarter — not alone. Get
    //           clear lessons, real-life practice, and a simple way to see how
    //           ready you are for citizenship, work, and daily life in the U.S.
    //         </p>
    //       </div>

    //       <div className="mt-6 flex flex-wrap gap-3">
    //         <button className="btn-primary bg-white !text-[var(--color-navy)]">
    //           Create a Free Account
    //         </button>
    //         <button className="btn-secondary border-white text-white bg-transparent">
    //           Explore Citizenship Prep
    //         </button>
    //       </div>
    //     </div>
    //   </div>
      
    // </section>