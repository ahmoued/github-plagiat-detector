import { Github, Twitter, Mail } from "lucide-react";

const Footer = () => {
  return (
    <footer className="border-t border-border bg-card/50 backdrop-blur-sm">
      <div className="container mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="flex flex-col md:flex-row items-center justify-between gap-4">
          <div className="text-center md:text-left">
            <p className="text-sm text-muted-foreground">
              © 2025 CodeGuard. Built for developers, by developers.
            </p>
          </div>
          
          <div className="flex items-center gap-4">
            <a
              href="https://github.com"
              target="_blank"
              rel="noopener noreferrer"
              className="p-2 rounded-lg bg-secondary/50 hover:bg-secondary transition-colors"
            >
              <Github className="w-5 h-5 text-foreground" />
            </a>
            <a
              href="https://twitter.com"
              target="_blank"
              rel="noopener noreferrer"
              className="p-2 rounded-lg bg-secondary/50 hover:bg-secondary transition-colors"
            >
              <Twitter className="w-5 h-5 text-foreground" />
            </a>
            <a
              href="mailto:contact@codeguard.dev"
              className="p-2 rounded-lg bg-secondary/50 hover:bg-secondary transition-colors"
            >
              <Mail className="w-5 h-5 text-foreground" />
            </a>
          </div>
        </div>
      </div>
    </footer>
  );
};

export default Footer;
