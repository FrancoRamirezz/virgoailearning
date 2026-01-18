"use client";

import Link from "next/link";
import Image from "next/image";
import { BookOpen } from "lucide-react";


export default function NavShell({ children, logoHref }: { children: React.ReactNode, logoHref: string }) {
    return (
        <nav className="bg-white/90 backdrop-blur-md shadow-lg border-b border-red-200 sticky top-0 z-50">
            <div className="px-4 mx-auto max-w-7xl">
                <div className="flex items-center justify-between h-16">
                    
                    <Link href={logoHref} className="flex items-center space-x-1">
                        <Image src="/logo.png" alt="VirgoLearning Logo" width={70} height={70} />
                        {/* <div className="w-10 h-10 bg-gradient-to-r from-red-600 to-blue-600 rounded-xl flex items-center justify-center">
                            
                        </div> */}
                        <span className="text-xl font-bold text-gray-900">
                        VirgoLearning
                        </span>
                    </Link>

                    {children}
                </div>
            </div>
        </nav>
    )
}    