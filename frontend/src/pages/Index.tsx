import { Link } from "react-router-dom";
import { Button } from "@/components/button";

const Index = () => {
  return (
    <div className="min-h-screen flex flex-col items-center justify-center p-4">
      <div className="max-w-2xl w-full text-center space-y-8">
        <div className="space-y-4">
          <h1 className="text-4xl sm:text-5xl font-bold">Code Checker</h1>

          <p className="text-lg text-muted-foreground">
            Check if your code is similar to other repositories
          </p>
        </div>

        <Link to="/checker">
          <Button size="lg" className="text-lg px-8">
            Start Checking
          </Button>
        </Link>
      </div>
    </div>
  );
};

export default Index;
