package io.jenkins.plugins.provenancegenerator;

import com.sun.jna.Structure;
import edu.umd.cs.findbugs.annotations.SuppressFBWarnings;

import java.nio.charset.StandardCharsets;
import java.util.Arrays;
import java.util.List;

@SuppressFBWarnings
public class GoString extends Structure implements Structure.ByValue {
    public String value;
    public long n;

    public GoString() {
    }

    @Override
    protected List<String> getFieldOrder() {
        return Arrays.asList("value", "n");
    }

    public static class ByReference extends GoString implements Structure.ByReference {
        public ByReference() {
        }

        public ByReference(String value) {
            super(value);
        }
    }

    public GoString(String value) {
        this.value = value == null ? "" : value;
        this.n = this.value.getBytes(StandardCharsets.UTF_8).length;
    }
}
