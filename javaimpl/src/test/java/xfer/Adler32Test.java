package xfer;

import junit.framework.Test;
import junit.framework.TestCase;
import junit.framework.TestSuite;

import java.io.UnsupportedEncodingException;

/**
 * Unit test for simple App.
 */
public class Adler32Test
    extends TestCase
{
   String small="hello this";
   String big="Apache is multithreaded: it spawns a thread per request (or process, it depends on the conf). \n" +
           "You can see how that overhead eats up memory as the number of concurrent connections increases \n" +
           "and more threads are needed to serve multiple simulataneous clients. Nginx and Node.js \n" +
           "are not multithreaded, because threads and processes carry a heavy memory cost. They are single-threaded, \n" +
           "but event-based. This eliminates \n" +
           "the overhead created by thousands of threads/processes by handling many connections in a single thread.";
    /**
     * Create the test case
     *
     * @param testName name of the test case
     */
    public Adler32Test(String testName)
    {
        super( testName );
    }

    /**
     * @return the suite of tests being tested
     */
    public static Test suite()
    {
        return new TestSuite( Adler32Test.class );
    }

    /**
     * Rigourous Test :-)
     */
    public void testApp()
    {
        Hashes app =new Hashes();
        byte[] stg=small.getBytes();
        long val=app.adler32(stg, 0, stg.length);
        System.out.println("testApp.val = " + val);
        java.util.zip.Adler32 jdkAdl=new java.util.zip.Adler32();
        jdkAdl.update(stg,0,stg.length);
        long valFromJdk=jdkAdl.getValue();
        System.out.println("testApp.valFromJdk = " + valFromJdk);
        assertEquals(valFromJdk,val);

    }


       public void testRoll() throws UnsupportedEncodingException {
        Hashes app =new Hashes();
        byte[] stg=big.getBytes("UTF-8");
        System.out.println("stg = " + stg.length);
        int mid=75;
        int st=17;
        app.adler32(stg, st, mid);
        while((st+mid<stg.length-1)) {
            System.out.println("st = " + st);
            long valR = app.roll(stg[st+mid]);
            System.out.println("testRoll.valR = " + valR);
            Hashes app2 =new Hashes();
            long val = app2.adler32(stg, st + 1, mid);
            System.out.println("testRoll.val = " + val);
            java.util.zip.Adler32 jdkAdl = new java.util.zip.Adler32();
            jdkAdl.update(stg, st+1, mid);
            long valFromJdk = jdkAdl.getValue();
            System.out.println("testRoll.valFromJdk = " + valFromJdk);
            if(val != valR){
                System.out.println(String.format("st = %s, len=%s,",new String[]{st+"",mid+""}));
            }
            assertEquals(valR, val);
            assertEquals(valR, valFromJdk);
            st++;
       }

    }

}
