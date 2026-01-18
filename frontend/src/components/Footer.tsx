import { BookOpen } from "lucide-react";
export default function Footer() {
  return (
    <footer className="bg-(--color-navy) text-white">
            <div className="px-4 py-10 mx-auto max-w-7xl lg:px-8">
              <div className="grid gap-8 md:grid-cols-2 lg:grid-cols-4">
                <div>
                  <div className="flex items-center space-x-3 mb-6">
                    <div className="w-10 h-10 bg-gradient-to-r from-red-600 to-blue-600 rounded-xl flex items-center justify-center">
                      <BookOpen className="w-6 h-6 text-white" />
                    </div>
                    <span className="text-xl font-bold text-(--color-red)">VirgoLearning</span>
                  </div>
                  <p className="text-gray-400 leading-relaxed">
                    Helping immigrants prepare for U.S. citizenship with clear lessons, real practice, and confidence-building tools.
                  </p>
                </div>
                
                <div>
                  <h3 className="text-lg font-semibold mb-4 text-(--color-red)">Practice & Learning</h3>
                  <ul className="space-y-2 text-(--color-red)">
                    <li><a href="#" >Citizenship Test Prep</a></li>
                    <li><a href="#" >English for Citizenship</a></li>
                    <li><a href="#" >Practice Tests</a></li>
                  </ul>
                </div>
                
                <div>
                  <h3 className="text-lg font-semibold mb-4 text-(--color-red)">Support</h3>
                  <ul className="space-y-2 text-gray-400">
                    <li><a href="#" >Help Center</a></li>
                    <li><a href="#" >Contact Us</a></li>
                    <li><a href="#" >FAQ</a></li>
                  </ul>
                </div>
                
                <div>
                  <h3 className="text-lg font-semibold mb-4 text-(--color-red)">Company/Legal</h3>
                  <ul className="space-y-2 text-gray-400">
                    <li><a href="#" >About Us</a></li>
                    <li><a href="#" >Privacy Policy</a></li>
                    <li><a href="#" >Terms of Service</a></li>
                  </ul>
                </div>
              </div>
              
              <div className="pt-8 mt-8 border-t border-gray-800 text-center text-gray-400">
                <p>&copy; 2025 CitizenshipPrep. All rights reserved.</p>
              </div>
            </div>
          </footer>
  )
}