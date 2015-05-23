package xfer;


import java.io.BufferedInputStream;
import java.io.FileInputStream;
import java.io.InputStream;
import java.util.ArrayList;
import java.util.List;

/**
 * Created by singhmo on 4/20/2015.
 */
public class Diff {

    public List<Object> diffData = new ArrayList<Object>();

    public static class ByteRange {
        int start;
        int end;

        ByteRange(int start, int end) {
            this.start = start;
            this.end = end;
        }

        @Override
        public boolean equals(Object o) {
            if (this == o) return true;
            if (o == null || getClass() != o.getClass()) return false;

            ByteRange range = (ByteRange) o;

            if (end != range.end) return false;
            if (start != range.start) return false;

            return true;
        }

        @Override
        public int hashCode() {
            int result = start;
            result = 31 * result + end;
            return result;
        }

        @Override
        public String toString() {
            return "ByteRange{" +
                    "start=" + start +
                    ", end=" + end +
                    '}';
        }
    }


}
