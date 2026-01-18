"use client";

export default function MobileMenu({open,children,}: {open: boolean; children: React.ReactNode;}) {
  if (!open) return null;

  return (
    <div className="md:hidden py-4 border-t border-red-200 bg-white/95 backdrop-blur-sm">
      <div className="flex flex-col space-y-1">{children}</div>
    </div>
  );
}
