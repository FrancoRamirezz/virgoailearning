"use client";

import { useState } from "react";
import { Menu, X, BookOpen, User, Info } from "lucide-react";
import NavShell from "./NavShell";
import NavLink from "./NavLink";
import MobileMenu from "./MobileMenu";

export default function AppNavbar({ role }: { role?: string }) {
  const [open, setOpen] = useState(false);

  return (
    <NavShell logoHref="/dashboard">
      <div className="hidden md:flex items-center space-x-2">
        <NavLink href="/dashboard" label="Dashboard" icon={BookOpen} />
        <NavLink href="/courses" label="Courses" icon={BookOpen} />
        <NavLink href="/contact" label="Contact" icon={BookOpen} />
        <NavLink href="/" label="Landing" icon={Info} />
      </div>

      <div className="hidden md:flex items-center space-x-3 bg-gray-50 px-3 py-2 rounded-lg">
        <div className="w-8 h-8 bg-gradient-to-r from-red-600 to-blue-600 rounded-full flex items-center justify-center">
                              <User className="w-5 h-5 text-white" />
                            </div>
        <span>{role ?? "Student"}</span>
      </div>

      <button onClick={() => setOpen(!open)} className="md:hidden">
        {open ? <X /> : <Menu />}
      </button>

      <MobileMenu open={open}>
        <NavLink href="/dashboard" label="Dashboard" icon={BookOpen} />
        <NavLink href="/civics" label="Civics" icon={BookOpen} />
        <NavLink href="/readiness" label="Readiness" icon={BookOpen} />
      </MobileMenu>
    </NavShell>
  );
}
