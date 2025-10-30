import { Link } from "react-router-dom";
import { Shield, Search, BarChart3, Zap, Github, CheckCircle } from "lucide-react";
import { Button } from "../components/button";
import Navbar from "../components/Navbar";
import Footer from "../components/Footer";
import { Card, CardContent } from "../components/card";

const Index = () => {
  const features = [
    {
      icon: Search,
      title: "Deep Code Analysis",
      description: "Advanced algorithms scan repositories for code similarities and patterns"
    },
    {
      icon: BarChart3,
      title: "Detailed Metrics",
      description: "Get comprehensive similarity scores and visual comparisons"
    },
    {
      icon: Zap,
      title: "Lightning Fast",
      description: "Results in seconds, not minutes. Optimized for speed and accuracy"
    },
    {
      icon: Github,
      title: "GitHub Integration",
      description: "Seamlessly works with any public GitHub repository"
    }
  ];

  const benefits = [
    "Protect your original code",
    "Verify project authenticity",
    "Educational integrity checks",
    "Open source compliance"
  ];

  return (
    <div className="min-h-screen flex flex-col">
      <Navbar />
      
      {/* Hero Section */}
      <section className="relative overflow-hidden pt-24 pb-32 px-4 sm:px-6 lg:px-8">
        <div className="absolute inset-0 bg-gradient-radial from-primary/10 via-transparent to-transparent opacity-50" />
        
        <div className="container mx-auto max-w-6xl relative">
          <div className="text-center animate-fade-in">
            <div className="inline-block mb-6 p-3 rounded-2xl bg-primary/10 backdrop-blur-sm">
              <Shield className="w-12 h-12 text-primary animate-glow" />
            </div>
            
            <h1 className="text-5xl sm:text-6xl lg:text-7xl font-bold mb-6 leading-tight">
              Detect Code Plagiarism
              <span className="block bg-gradient-to-r from-primary via-accent to-primary bg-clip-text text-transparent">
                Protect Your Work
              </span>
            </h1>
            
            <p className="text-xl text-muted-foreground max-w-3xl mx-auto mb-12 leading-relaxed">
              Advanced plagiarism detection for GitHub repositories. Analyze code similarities, 
              identify potential duplicates, and ensure originality with our powerful scanning engine.
            </p>
            
            <div className="flex flex-col sm:flex-row items-center justify-center gap-4 animate-fade-in-up">
              <Link to="/checker">
                <Button
                  size="lg"
                  className="h-14 px-8 text-lg bg-primary hover:bg-primary/90 text-primary-foreground font-semibold shadow-lg shadow-primary/30 hover:shadow-xl hover:shadow-primary/40 transition-all"
                >
                  <Search className="w-5 h-5 mr-2" />
                  Start Checking Now
                </Button>
              </Link>
              
              <Button
                size="lg"
                variant="outline"
                className="h-14 px-8 text-lg border-border hover:border-primary hover:bg-secondary transition-all"
                asChild
              >
                <a href="#features">
                  Learn More
                </a>
              </Button>
            </div>
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section id="features" className="py-20 px-4 sm:px-6 lg:px-8 bg-card/30">
        <div className="container mx-auto max-w-6xl">
          <div className="text-center mb-16 animate-fade-in">
            <h2 className="text-4xl font-bold mb-4 text-foreground">
              Powerful Features
            </h2>
            <p className="text-lg text-muted-foreground max-w-2xl mx-auto">
              Everything you need to verify code originality and maintain integrity
            </p>
          </div>

          <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-8">
            {features.map((feature, index) => (
              <Card
                key={index}
                className="group hover:shadow-lg hover:shadow-primary/10 transition-all duration-300 hover:-translate-y-2 bg-gradient-to-br from-card to-card/50 border-border/50 animate-fade-in-up"
                style={{ animationDelay: `${index * 0.1}s` }}
              >
                <CardContent className="p-6">
                  <div className="p-3 rounded-xl bg-primary/10 inline-block mb-4 group-hover:bg-primary/20 transition-colors">
                    <feature.icon className="w-6 h-6 text-primary" />
                  </div>
                  <h3 className="text-xl font-semibold mb-2 text-foreground">
                    {feature.title}
                  </h3>
                  <p className="text-muted-foreground">
                    {feature.description}
                  </p>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>
      </section>

      {/* Benefits Section */}
      <section className="py-20 px-4 sm:px-6 lg:px-8">
        <div className="container mx-auto max-w-6xl">
          <div className="grid lg:grid-cols-2 gap-12 items-center">
            <div className="animate-fade-in">
              <h2 className="text-4xl font-bold mb-6 text-foreground">
                Why Choose CodeGuard?
              </h2>
              <p className="text-lg text-muted-foreground mb-8">
                Built for developers, educators, and organizations who value code integrity 
                and want to ensure originality in their projects.
              </p>
              
              <div className="space-y-4">
                {benefits.map((benefit, index) => (
                  <div
                    key={index}
                    className="flex items-center gap-3 p-4 rounded-lg bg-secondary/30 hover:bg-secondary/50 transition-colors animate-fade-in-up"
                    style={{ animationDelay: `${index * 0.1}s` }}
                  >
                    <CheckCircle className="w-5 h-5 text-primary flex-shrink-0" />
                    <span className="text-foreground font-medium">{benefit}</span>
                  </div>
                ))}
              </div>
            </div>

            <div className="relative animate-scale-in">
              <div className="absolute inset-0 bg-gradient-to-r from-primary/20 to-accent/20 rounded-3xl blur-3xl" />
              <Card className="relative bg-gradient-to-br from-card to-card/50 border-border/50 overflow-hidden">
                <CardContent className="p-8">
                  <div className="space-y-6">
                    <div className="flex items-center justify-between p-4 rounded-lg bg-secondary/50">
                      <span className="text-sm font-medium text-muted-foreground">Similarity Score</span>
                      <span className="text-2xl font-bold text-primary">87%</span>
                    </div>
                    <div className="space-y-3">
                      <div className="h-3 bg-secondary rounded-full overflow-hidden">
                        <div className="h-full bg-gradient-to-r from-primary to-accent w-[87%] rounded-full animate-fade-in" />
                      </div>
                      <div className="h-3 bg-secondary rounded-full overflow-hidden">
                        <div className="h-full bg-gradient-to-r from-primary to-accent w-[65%] rounded-full animate-fade-in" style={{ animationDelay: "0.2s" }} />
                      </div>
                      <div className="h-3 bg-secondary rounded-full overflow-hidden">
                        <div className="h-full bg-gradient-to-r from-primary to-accent w-[42%] rounded-full animate-fade-in" style={{ animationDelay: "0.4s" }} />
                      </div>
                    </div>
                    <div className="pt-4 border-t border-border/50">
                      <p className="text-sm text-muted-foreground">
                        Real-time analysis with detailed metrics and insights
                      </p>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </div>
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-20 px-4 sm:px-6 lg:px-8 bg-gradient-to-b from-card/30 to-transparent">
        <div className="container mx-auto max-w-4xl text-center animate-fade-in">
          <h2 className="text-4xl font-bold mb-6 text-foreground">
            Ready to Get Started?
          </h2>
          <p className="text-lg text-muted-foreground mb-8 max-w-2xl mx-auto">
            Start checking your repositories for plagiarism today. Fast, accurate, and easy to use.
          </p>
          <Link to="/checker">
            <Button
              size="lg"
              className="h-14 px-8 text-lg bg-primary hover:bg-primary/90 text-primary-foreground font-semibold shadow-lg shadow-primary/30 hover:shadow-xl hover:shadow-primary/40 transition-all"
            >
              <Search className="w-5 h-5 mr-2" />
              Check Your First Repository
            </Button>
          </Link>
        </div>
      </section>

      <Footer />
    </div>
  );
};

export default Index;
