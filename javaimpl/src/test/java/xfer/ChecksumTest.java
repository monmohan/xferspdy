package xfer;

import junit.framework.Test;
import junit.framework.TestCase;
import junit.framework.TestSuite;

import java.io.*;
import java.nio.channels.FileChannel;
import java.util.List;

/**
 * Created by singhmo on 4/21/2015.
 */
public class ChecksumTest extends TestCase{
    public static Test suite()
    {
        return new TestSuite( ChecksumTest.class );
    }

    public ChecksumTest(String tName) {
        super(tName);
    }

    public void testcsgen() throws Exception {
        File file=new File("D:\\work\\projects\\xferspdy\\temp\\in.txt");
        File ofile=new File("D:\\work\\projects\\xferspdy\\temp\\out.txt");
        try {
            String test="The quick brown fox jumps over the lazy dog";
            StringBuilder b=new StringBuilder();
            for(int i=0;i<100;i++){
                b.append(test);
            }
            FileWriter fc=new FileWriter(file);
            fc.write(b.toString());
            fc.close();
            Checksum cs=Serializer.toChecksum(new DataSource.FileDataSource(file),64);
            System.out.println("cs = " + cs.fingerprint);

            Serializer.fromChecksum(cs,new FileOutputStream(ofile));
            FileReader fin=new FileReader(ofile);
            char[] buf=new char[b.length()];
            fin.read(buf);
            assertEquals(new String(buf),b.toString());
        } finally {
            //file.delete();
            //ofile.delete();
        }


    }

    public void testRmLastBlock() throws Exception {
        byte[] b=new byte[1000];
        for (int i = 0; i < b.length; i++) {
            b[i] = (byte)(i%128);
        }
        DataSource ds=new DataSource.MemDataSource(b);
        Checksum in= Serializer.toChecksum(ds, 100);
        System.out.println("in = " + in.blocks);
        Diff diff=Serializer.diff(in,new DataSource.MemDataSource(b,0,900));
        System.out.println("diff.diffData = " + diff.diffData);
    }

    public void testChangeFirstByte() throws Exception {
        byte[] b=new byte[100];
        for (int i = 0; i < b.length; i++) {
            b[i] = (byte)(i%128);
        }

        DataSource ds=new DataSource.MemDataSource(b);
        Checksum in= Serializer.toChecksum(ds, 10);
        System.out.println("in = " + in.blocks);

        b[0]=31;
        Diff diff=Serializer.diff(in,new DataSource.MemDataSource(b));

        System.out.println("diff.diffData = " + diff.diffData);
    }


}
