import Link from "next/link";
import { usePathname } from "next/navigation";

export default function NavLink({href,icon:Icon,label}:{href:string,icon:React.ComponentType<any>,label:string}) {
    const pathname = usePathname();
    const isActive = pathname === href;
    return (
        <Link
            href={href}
            className={`flex items-center space-x-2 px-3 py-2 rounded-lg transition-colors duration-200 ${
                isActive
                ? "text-red-600 bg-red-50 shadow-sm"
                : "text-gray-600 hover:text-gray-900 hover:bg-gray-50"
            }`}
            >
            <Icon className="w-4 h-4" />
            <span className="font-medium">{label}</span>
        </Link>
    )
}