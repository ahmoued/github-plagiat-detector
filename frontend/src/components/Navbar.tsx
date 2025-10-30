import { Link, useLocation } from "react-router-dom";
import { Shield } from "lucide-react";

const Navbar = () => {
  const location = useLocation();

  return (
    <nav className="fixed top-0 w-full z-50 bg-background/80 backdrop-blur-md border-b border-border">
      <div className="container mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          <Link to="/" className="flex items-center gap-2 group">
            <div className="p-2 rounded-lg bg-primary/10 group-hover:bg-primary/20 transition-colors">
              <Shield className="w-6 h-6 text-primary" />
            </div>
            <span className="text-xl font-bold bg-gradient-to-r from-primary to-accent bg-clip-text text-transparent">
              CodeGuard
            </span>
          </Link>
          
          <div className="flex items-center gap-6">
            <Link
              to="/"
              className={`text-sm font-medium transition-colors hover:text-primary ${
                location.pathname === "/" ? "text-primary" : "text-muted-foreground"
              }`}
            >
              Home
            </Link>
            <Link
              to="/checker"
              className={`text-sm font-medium transition-colors hover:text-primary ${
                location.pathname === "/checker" ? "text-primary" : "text-muted-foreground"
              }`}
            >
              Checker
            </Link>
          </div>
        </div>
      </div>
    </nav>
  );
};

export default Navbar;
