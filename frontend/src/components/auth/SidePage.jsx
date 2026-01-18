
import Image from "next/image"
export default function SidePage() {
  return (
      <div className="hidden lg:flex lg:w-1/2 bg-gradient-to-bl from-red-600 to-blue-600 p-12 items-center justify-center relative overflow-hidden">
        {/* Decorative Elements */}
        <div className="absolute top-0 right-0 w-96 h-96 bg-white/5 rounded-full -translate-y-48 translate-x-48"></div>
        <div className="absolute bottom-0 left-0 w-80 h-80 bg-white/5 rounded-full translate-y-40 -translate-x-40"></div>
        
        <div className="relative z-10 max-w-lg text-white">
          <div className="mb-8">
            <div className="inline-flex items-center gap-2 px-4 py-2 bg-white/10 backdrop-blur-sm rounded-full mb-6">
              <div className="w-2 h-2 bg-green-400 rounded-full animate-pulse"></div>
              <span className="text-sm font-medium">Built for future U.S. citizens</span>
            </div>
            
            <h2 className="text-4xl font-bold mb-4">
              Get Ready for Your U.S. Citizenship Test
            </h2>
            <p className="text-lg text-red-100 mb-8">
              CitizenshipPrep helps immigrants study civics, practice English, and track their readiness 
              so they can walk into their citizenship interview feeling confident and prepared.
            </p>
          </div>

          {/* Stats */}
          <div className="grid grid-cols-3 gap-6">
            <div className="bg-white/10 backdrop-blur-sm rounded-xl p-4">
              <div className="text-3xl font-bold mb-1">100</div>
              <div className="text-sm text-red-100">Official Civics Questions</div>
            </div>
            <div className="bg-white/10 backdrop-blur-sm rounded-xl p-4">
              <div className="text-3xl font-bold mb-1">AI</div>
              <div className="text-sm text-red-100">Personalized Practice</div>
            </div>
            <div className="bg-white/10 backdrop-blur-sm rounded-xl p-4">
              <div className="text-3xl font-bold mb-1">1</div>
              <div className="text-sm text-red-100">Clear Readiness Score</div>
            </div>
          </div>

          {/* Testimonial */}
          <div className="mt-12 bg-white/10 backdrop-blur-sm rounded-xl p-6">
            <p className="text-white/90 italic mb-4">
              "Our goal is simple: help every immigrant know exactly how ready they are for their 
              U.S. citizenship interview and what to study next."
            </p>
            <div className="flex items-center gap-3">
              <Image src="/logo.png" alt="VirgoLearning Logo" width={50} height={70} />
              {/* <div className="w-10 h-10 bg-white/20 rounded-full"></div>*/}
              <div> 
                <div className="font-semibold">The CitizenshipPrep Team</div>
                <div className="text-sm text-red-200">Building tools for your journey</div>
              </div>
            </div>
          </div>
        </div>
      </div>
  );
}