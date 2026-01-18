import { BookOpen } from "lucide-react";
import Image from "next/image";

export default function PublicFooter() {
  return (
    <footer className="bg-(--color-blue-600) text-white">
      <div className="px-4 py-12 mx-auto max-w-7xl lg:px-8">
        <div className="grid gap-10 md:grid-cols-2 lg:grid-cols-4">
          
          {/* Brand + Mission */}
          <div>
            <div className="flex items-center space-x-3 mb-6">
              
              <div className="w-15 h-15 bg-white rounded-3xl flex items-center justify-center">
                <Image src="/logo.png" alt="VirgoLearning Logo" width={70} height={70} />
              </div>
              <span className="text-xl font-bold text-white">
                VirgoLearning
              </span>
            </div>

            <p className="text-gray-400 leading-relaxed">
              Helping immigrants prepare for U.S. citizenship with clear lessons,
              real practice, and confidence-building tools.
            </p>
          </div>

          {/* Practice & Learning */}
          <div>
            <h3 className="text-lg font-semibold mb-4 text-white">
              Practice & Learning
            </h3>
            <ul className="space-y-2 text-gray-400">
              <li><a href="/practice/civics" className="hover:text-white">Citizenship Test Prep</a></li>
              <li><a href="/practice/english" className="hover:text-white">English for Citizenship</a></li>
              <li><a href="/practice/tests" className="hover:text-white">Practice Tests</a></li>
            </ul>
          </div>

          {/* Support */}
          <div>
            <h3 className="text-lg font-semibold mb-4 text-white">
              Support
            </h3>
            <ul className="space-y-2 text-gray-400">
              <li><a href="/help" className="hover:text-white">Help Center</a></li>
              <li><a href="/faq" className="hover:text-white">FAQ</a></li>
              <li><a href="/contact" className="hover:text-white">Contact Us</a></li>
            </ul>
          </div>

          {/* Company / Legal */}
          <div>
            <h3 className="text-lg font-semibold mb-4 text-white">
              Company & Legal
            </h3>
            <ul className="space-y-2 text-gray-300">
              <li><a href="/about" className="hover:text-white">About Us</a></li>
              <li><a href="/privacy" className="hover:text-white">Privacy Policy</a></li>
              <li><a href="/terms" className="hover:text-white">Terms of Service</a></li>
            </ul>
          </div>
        </div>

        {/* Bottom */}
        <div className="pt-8 mt-10 border-t border-red-800 text-center text-gray-300 text-sm">
          <p>&copy; 2025 VirgoLearning. All rights reserved.</p>
          <p className="mt-2">
            VirgoLearning provides educational tools only and does not offer legal advice.
          </p>
        </div>
      </div>
    </footer>
  );
}
