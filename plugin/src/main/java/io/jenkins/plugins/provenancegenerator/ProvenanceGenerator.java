package io.jenkins.plugins.provenancegenerator;


import com.sun.jna.Library;
import com.sun.jna.Native;
import hudson.Extension;
import hudson.FilePath;
import hudson.Launcher;
import hudson.model.AbstractProject;
import hudson.model.Run;
import hudson.model.TaskListener;
import hudson.tasks.BuildStepDescriptor;
import hudson.tasks.Builder;
import jenkins.tasks.SimpleBuildStep;
import org.jenkinsci.Symbol;
import org.kohsuke.stapler.DataBoundConstructor;

import java.io.*;

public class ProvenanceGenerator extends Builder implements SimpleBuildStep {
    private final String artifact;
    //private final libProvenanceGeneartor INSTANCE;
    private final String output;
    private final String envpath;

    @DataBoundConstructor
    public ProvenanceGenerator(String artifact, String output, String envpath) {
        this.artifact = artifact;
        this.output = output;
        this.envpath = envpath;
    }

    public String getArtifact() {
        return artifact;
    }
    public String getOutput() {
        return output;
    }
    public String getEnvpath() {
        return envpath;
    }

    @Override
    public void perform(Run build,
                        FilePath workspace,
                        Launcher launcher,
                        TaskListener listener) throws InterruptedException, IOException {
        generateProvenace(artifact,output,envpath);
    }

    public void generateProvenace(String artifact, String output, String envpath){
        InputStream is;//classloader.getResourceAsStream("/io/jenkins/plugins/sample/provenance.dylib");
        File file = null;
        String osName = System.getProperty("os.name");
        if(osName.startsWith("Mac OS X")){
            try {
                is = ProvenanceGenerator.class.getResourceAsStream("provenance.dylib");
                file = File.createTempFile("lib", ".dylib");
            } catch (IOException e) {
                throw new RuntimeException(e);
            }
        } else if (osName.startsWith("Windows")) {
            try {
                is = ProvenanceGenerator.class.getResourceAsStream("provenance.dll");
                file = File.createTempFile("lib", ".dll");
            } catch (IOException e) {
                throw new RuntimeException(e);
            }
        }
        else{ //linux
            try {
                is = ProvenanceGenerator.class.getResourceAsStream("provenance.so");
                file = File.createTempFile("lib", ".so");
            } catch (IOException e) {
                throw new RuntimeException(e);
            }
        }

        OutputStream os = null;
        try {
            os = new FileOutputStream(file);
        } catch (FileNotFoundException e) {
            throw new RuntimeException(e);
        }
        byte[] buffer = new byte[1024];
        int length;
        while (true) {
            try {
                if (!((length = is.read(buffer)) != -1)) break;
            } catch (IOException e) {
                throw new RuntimeException(e);
            }
            try {
                os.write(buffer, 0, length);
            } catch (IOException e) {
                throw new RuntimeException(e);
            }
        }
        try {
            is.close();
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
        try {
            os.close();
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
        libProvenanceGeneartor INSTANCE = Native.load(file.getAbsolutePath(), libProvenanceGeneartor.class);
        INSTANCE.GenerateSLSA(new GoString(artifact),new GoString(output), new GoString(envpath));
        file.deleteOnExit();
    }

    public interface libProvenanceGeneartor extends Library {
        void GenerateSLSA(GoString artifact, GoString outpath, GoString envpath);
    }

    @Symbol("provenance-generator")
    @Extension
    public static class DescriptorImpl extends BuildStepDescriptor<Builder> {

        @Override
        public String getDisplayName() {
            return "Provenance Generator";
        }

        @Override
        public boolean isApplicable(Class<? extends AbstractProject> t) {
            return true;
        }
    }
}
