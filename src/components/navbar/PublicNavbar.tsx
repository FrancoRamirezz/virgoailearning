"use client";

import { useState } from "react";
import { Menu, X, Info, DollarSign, Phone, BookOpen } from "lucide-react";
import NavShell from "./NavShell";
import NavLink from "./NavLink";
import MobileMenu from "./MobileMenu";
import Link from "next/link";

export default function PublicNavbar() {
  const [open, setOpen] = useState(false);

  return (
    <NavShell logoHref="/">
      {/* Desktop */}
      <div className="hidden md:flex items-center space-x-2">
        <NavLink href="/" label="Home" icon={BookOpen} />
        <NavLink href="/about" label="About" icon={Info} />
        <NavLink href="/pricing" label="Pricing" icon={DollarSign} />
        <NavLink href="/dashboard" label="dashboard" icon={Info} />
      </div>

      <div className="hidden md:flex items-center space-x-4">
        <Link href="/login" className="text-black hover:text-(--color-red)">Log In</Link>
        <Link href="/signup" className="bg-gradient-to-r from-red-600 to-blue-600 text-white px-4 py-2 rounded-lg">
          Sign Up
        </Link>
      </div>

      {/* Mobile */}
      <button onClick={() => setOpen(!open)} className="md:hidden">
        {open ? <X /> : <Menu />}
      </button>

      <MobileMenu open={open}>
        <NavLink href="/" label="Home" icon={Info} />
        <NavLink href="/about" label="About" icon={Info} />
        <NavLink href="/pricing" label="Pricing" icon={DollarSign} />
      </MobileMenu>
    </NavShell>
  );
}
